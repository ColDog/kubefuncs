package kubefuncs

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/golang/protobuf/proto"
	nsq "github.com/nsqio/go-nsq"
)

// Run sets up an opinionated client.
func Run(handler Handler, options ...Option) error {
	done := make(chan struct{}, 1)

	// TODO: Close when receive exit signal.

	c, err := NewClient(options...)
	if err != nil {
		return err
	}
	c.ListenHealthz()
	c.OnDefault(handler)

	<-done
	c.Close()
	return nil
}

// NewClient initializes a new client with the provided options. An error will
// be returned if producers and consumers cannot be setup.
func NewClient(options ...Option) (*Client, error) {
	opts := defaults()
	for _, o := range options {
		o(opts)
	}
	c, err := newNsqClient(opts.nsqdURL, opts.lookupdURL)
	if err != nil {
		return nil, err
	}

	h := &Client{
		opts:      *opts,
		nsq:       c,
		responses: map[string]chan *Event{},
	}
	if opts.rpc {
		err = h.setupReturn(opts)
		if err != nil {
			return nil, err
		}
	}
	log.Printf("client initialized: %+v", *opts)
	return h, nil
}

// Client is the core client library to interact with the function architecture.
type Client struct {
	opts

	nsq         nsqClient
	returnTopic string
	lock        sync.RWMutex

	responses map[string]chan *Event
}

func (h *Client) setupReturn(opts *opts) error {
	h.returnTopic = "rpc-" + opts.clientID + "#ephemeral"
	if err := h.Emit(
		context.Background(), &Event{Id: "ping", Topic: h.returnTopic}); err != nil {
		return err
	}

	if err := h.nsq.subscribe(
		h.returnTopic,
		"default#ephemeral",
		nsq.HandlerFunc(h.handleResponse),
	); err != nil {
		return err
	}
	return nil
}

// ListenHealthz opens an http server to provide health checks. Returns after
// starting the server in a goroutine.
func (h *Client) ListenHealthz() {
	go func() {
		err := http.ListenAndServe(h.healthzAddr,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			}),
		)
		if err != nil {
			log.Printf("healthz err: %v", err)
		}
	}()
}

// Close all producers and consumers.
func (h *Client) Close() { h.nsq.close() }

// Emit will asyncronously publish a message and will not wait for any
// responses.
func (h *Client) Emit(ctx context.Context, ev *Event) error {
	data, err := proto.Marshal(ev)
	if err != nil {
		return err
	}
	return h.nsq.publish(ev.Topic, data)
}

// Call will publish a message and wait on a return queueu.
func (h *Client) Call(ctx context.Context, ev *Event) (*Message, error) {
	if !h.rpc {
		return nil, errors.New("call disabled: use WithCallEnabled() when initializing the client")
	}
	ev.Return = h.returnTopic
	if err := h.Emit(ctx, ev); err != nil {
		return nil, err
	}
	resp, err := h.wait(ctx, ev.Id)
	if err != nil {
		return nil, err
	}
	return &Message{Event: resp}, nil
}

// On will execute the provided handler for all messages along the provided
// queue and channel.
func (h *Client) On(topic, channel string, handler Handler) {
	// Publish a specific ping message to the topic, let's see if we get it.
	// This is a minor sanity check.
	h.Emit(context.Background(), &Event{Id: "ping", Topic: topic})

	h.nsq.subscribe(topic, channel, nsq.HandlerFunc(func(msg *nsq.Message) error {
		// Unmarshal the provided event. If this fails, this message should just
		// be removed as it is unprocessable.
		ev := &Event{}
		if err := proto.Unmarshal(msg.Body, ev); err != nil {
			log.Printf("failing marshal: %v", err)
			return nil // Return nil to not requeue.
		}

		// If the event id is a ping then we should continue.
		if ev.Id == "ping" {
			log.Println("ping received")
			return nil
		}

		// Wrap the event in the message struct.
		wrap := &Message{Event: ev}
		log.Printf("handling: %v", ev)

		// Handle the provided message. If we get an error message and we have
		// a return queue, we return the message and finish the response.
		if err := handler.Handle(wrap); err != nil {
			if ev.Return != "" {
				// Add a response.
				err := wrap.Respond(&Error{Error: err.Error()})
				log.Printf("handler error: %v, %+v", err, wrap.Response)
			} else {
				// If we've received a handler error, and no return queue, respond
				// with an error so this is re-queued.
				return err
			}
		}

		// If we processed successfully and the return queue is present, we
		// should respond.
		if ev.Return != "" {
			// If no response is present, we add an empty response.
			if wrap.Response == nil {
				if err := wrap.Respond(&Empty{}); err != nil {
					return err // This should never happen.
				}
			}

			// Marshal the response value.
			data, err := proto.Marshal(wrap.Response)
			if err != nil {
				return err
			}

			log.Printf("publishing return: %+v", wrap.Response)

			// Publish the response.
			return h.nsq.publish(ev.Return, data)
		}
		return nil
	}))
}

// OnDefault executes the handler for the default configured topic and channel.
func (h *Client) OnDefault(handler Handler) {
	h.On(h.topic, h.channel, handler)
}

func (h *Client) handleResponse(msg *nsq.Message) error {
	defer msg.Finish()

	r := &Event{}
	err := proto.Unmarshal(msg.Body, r)
	if err != nil {
		log.Printf("failed to decode: %v", err)
		return nil
	}

	if r.Id == "ping" {
		log.Printf("ping successful")
		return nil
	}

	h.lock.RLock()
	c, ok := h.responses[r.Id]
	h.lock.RUnlock()
	if !ok {
		log.Printf("response not found %s: %+v", r.Id, h.responses)
		return nil
	}

	select {
	case c <- r:
	default:
		log.Printf("deadlock")
	}
	return nil
}

func (h *Client) wait(ctx context.Context, id string) (*Event, error) {
	c := make(chan *Event, 1)
	h.lock.Lock()
	h.responses[id] = c
	h.lock.Unlock()
	defer func() {
		h.lock.Lock()
		delete(h.responses, id)
		h.lock.Unlock()
	}()

	select {
	case res := <-c:
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
