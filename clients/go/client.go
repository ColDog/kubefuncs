package kubefuncs

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/coldog/kubefuncs/clients/go/message"
	"github.com/golang/protobuf/proto"
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
	var nsq message.Client

	if opts.mockClient {
		nsq = message.NewMockClient()
	} else {
		c, err := message.NewNSQClient(opts.nsqdURL, opts.lookupdURL, opts.logger)
		if err != nil {
			return nil, err
		}
		nsq = c
	}

	h := &Client{
		opts:      *opts,
		nsq:       nsq,
		responses: map[string]chan *Event{},
	}
	if opts.rpc {
		if err := h.setupReturn(opts); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// Client is the core client library to interact with the function architecture.
type Client struct {
	opts

	nsq         message.Client
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

	if err := h.nsq.Subscribe(
		h.returnTopic, "default#ephemeral", h.handleResponse); err != nil {
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
			h.log(0, "healthz err: %v", err)
		}
	}()
}

// Close all producers and consumers.
func (h *Client) Close() { h.nsq.Close() }

// Emit will asyncronously publish a message and will not wait for any
// responses.
func (h *Client) Emit(ctx context.Context, ev *Event) error {
	data, err := proto.Marshal(ev)
	if err != nil {
		return err
	}
	return h.nsq.Publish(ev.Topic, data)
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
	h.nsq.Subscribe(topic, channel, func(body []byte) error {
		// Unmarshal the provided event. If this fails, this message should just
		// be removed as it is unprocessable.
		ev := &Event{}
		if err := proto.Unmarshal(body, ev); err != nil {
			h.log(4, "failing marshal: %v", err)
			return nil // Return nil to not requeue.
		}

		// If the event id is a ping then we should continue. This is to allow
		// for healthchecks.
		if ev.Id == "ping" {
			h.log(4, "ping received")
			return nil
		}

		// Wrap the event in the message struct.
		wrap := &Message{Event: ev}
		h.log(5, "handling: %v", ev)

		// Handle the provided message. If we get an error message and we have
		// a return queue, we return the message and finish the response.
		if err := handler.Handle(wrap); err != nil {
			if ev.Return != "" {
				// Add a response.
				err := wrap.Respond(&Error{Error: err.Error()})
				h.log(4, "handler error: %v, %+v", err, wrap.Response)
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

			h.log(4, "publishing return: %+v", wrap.Response)

			// Publish the response.
			return h.nsq.Publish(ev.Return, data)
		}
		return nil
	})
}

// OnDefault executes the handler for the default configured topic and channel.
func (h *Client) OnDefault(handler Handler) {
	h.On(h.topic, h.channel, handler)
}

func (h *Client) handleResponse(body []byte) error {
	r := &Event{}
	err := proto.Unmarshal(body, r)
	if err != nil {
		h.log(0, "failed to decode: %v", err)
		return nil
	}

	if r.Id == "ping" {
		h.log(4, "ping successful")
		return nil
	}

	h.lock.RLock()
	c, ok := h.responses[r.Id]
	h.lock.RUnlock()
	if !ok {
		h.log(4, "response not found %s: %+v", r.Id, h.responses)
		return nil
	}

	c <- r
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

func (h *Client) log(level uint, msg string, args ...interface{}) {
	if h.logger != nil && level >= h.logVerbosity {
		h.logger.Printf(msg, args...)
	}
}
