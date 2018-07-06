package kubefuncs

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
	nsq "github.com/nsqio/go-nsq"
)

func Run(handler Handler, options ...Option) error {
	c, err := NewClient(options...)
	if err != nil {
		return err
	}
	c.ListenHealthz()
	c.OnDefault(handler)
	c.Wait()
	return nil
}

func NewClient(options ...Option) (*Client, error) {
	opts := defaults()
	for _, o := range options {
		o(opts)
	}
	h := &Client{
		opts:      *opts,
		responses: map[string]chan *Event{},
	}
	if err := h.setupProducer(); err != nil {
		return nil, err
	}
	if opts.rpc {
		h.returnTopic = "rpc-" + opts.clientID + "#ephemeral"
		if err := h.Emit(
			context.Background(), &Event{Id: "ping", Topic: h.returnTopic}); err != nil {
			return nil, err
		}

		if err := h.subscribe(
			h.returnTopic,
			"default#ephemeral",
			nsq.HandlerFunc(h.handleResponse),
		); err != nil {
			return nil, err
		}
	}
	log.Printf("client initialized: %+v", *opts)
	return h, nil
}

type Client struct {
	opts

	returnTopic string
	lock        sync.RWMutex
	producer    *nsq.Producer
	consumers   []*nsq.Consumer

	responses map[string]chan *Event
}

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

func (h *Client) Wait() {
	for _, c := range h.consumers {
		<-c.StopChan
	}
}

func (h *Client) Close() {
	if h.producer != nil {
		h.producer.Stop()
	}
	for _, c := range h.consumers {
		c.Stop()
	}
}

func (h *Client) Emit(ctx context.Context, ev *Event) error {
	data, err := proto.Marshal(ev)
	if err != nil {
		return err
	}
	return h.producer.PublishAsync(ev.Topic, data, nil)
}

func (h *Client) Call(ctx context.Context, ev *Event) (*Message, error) {
	if !h.rpc {
		return nil, errors.New("call disabled: use WithCallEnabled() when initializing the client")
	}
	ev.Return = h.returnTopic
	if err := h.Emit(ctx, ev); err != nil {
		return nil, err
	}
	h.register(ev.Id)
	resp, err := h.wait(ctx, ev.Id)
	if err != nil {
		return nil, err
	}
	return &Message{Event: resp}, nil
}

func (h *Client) OnDefault(handler Handler) {
	h.On(h.topic, h.channel, handler)
}

func (h *Client) On(topic, channel string, handler Handler) {
	h.Emit(context.Background(), &Event{Id: "ping", Topic: topic})

	h.subscribe(topic, channel, nsq.HandlerFunc(func(msg *nsq.Message) error {
		ev := &Event{}
		if err := proto.Unmarshal(msg.Body, ev); err != nil {
			return err
		}
		wrap := &Message{Event: ev}
		log.Printf("handling: %v", ev)

		if ev.Id == "ping" {
			msg.Finish()
			return nil
		}

		if err := handler.Handle(wrap); err != nil {
			log.Printf("handler err: %v", err)

			if ev.Return != "" {
				data, err := proto.Marshal(&Error{
					Error: err.Error(),
				})
				if err != nil {
					return err
				}

				h.producer.PublishAsync(ev.Return, data, nil)
				msg.Finish()
				return nil
			}

			msg.Requeue(-1)
			return err
		}

		if ev.Return != "" {
			if wrap.Response == nil {
				returnVal, err := NewEvent(ev.Return, &Empty{})
				if err != nil {
					return err
				}
				wrap.Response = returnVal
			}
			wrap.Response.Id = wrap.Event.Id
			data, err := proto.Marshal(wrap.Response)
			if err != nil {
				log.Printf("marshal err: %v", ev)
				return err
			}
			if err := h.producer.Publish(ev.Return, data); err != nil {
				log.Printf("return err: %v", ev)
				return err
			}
			log.Println("published")
		}

		msg.Finish()
		return nil
	}))
}

func (h *Client) setupProducer() error {
	producer, err := nsq.NewProducer(h.nsqdURL, nsq.NewConfig())
	if err != nil {
		return err
	}
	producer.SetLogger(log.New(os.Stderr, "[nsq] ", log.LstdFlags), nsq.LogLevelDebug)
	h.producer = producer
	return nil
}

func (h *Client) subscribe(topic, channel string, handler nsq.Handler) error {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}
	consumer.SetLogger(log.New(os.Stderr, "[nsq] ", log.LstdFlags), nsq.LogLevelDebug)
	consumer.AddHandler(handler)
	h.consumers = append(h.consumers, consumer)
	err = consumer.ConnectToNSQLookupd(h.lookupdURL)
	if err != nil {
		return err
	}
	return nil
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

	c <- r
	return nil
}

func (h *Client) register(id string) {
	h.lock.Lock()
	h.responses[id] = make(chan *Event)
	h.lock.Unlock()
}

func (h *Client) wait(ctx context.Context, id string) (*Event, error) {
	h.lock.RLock()
	c, ok := h.responses[id]
	h.lock.RUnlock()
	if !ok {
		return nil, errors.New("not found")
	}
	select {
	case res := <-c:
		return res, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
