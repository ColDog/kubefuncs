package kubefuncs

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"
)

// Message wraps an event with a Respond and Payload unmarshal option.
type Message struct {
	*Event
	Response *Event
}

// NewEvent initializes a new event given a payload.
func NewEvent(topic string, payload proto.Message) (*Event, error) {
	if topic == "" {
		return nil, errors.New("topic is empty")
	}

	id := uuid.NewV4().String()
	a, err := ptypes.MarshalAny(payload)
	if err != nil {
		return nil, err
	}

	ev := &Event{
		Id:      id,
		Topic:   topic,
		Payload: a,
	}
	return ev, nil
}

// Payload will decode the payload into a protobuf.
func (m *Message) Payload(dst proto.Message) error {
	return ptypes.UnmarshalAny(m.Event.Payload, dst)
}

// Respond will respond to the message.
func (m *Message) Respond(payload proto.Message) error {
	if m.Return == "" {
		return fmt.Errorf("cannot respond to this message: %s", m.Id)
	}
	a, err := ptypes.MarshalAny(payload)
	if err != nil {
		return err
	}
	m.Response = &Event{
		Id:      m.Event.Id,
		Topic:   m.Event.Return,
		Payload: a,
	}
	return nil
}

// Handler represents the message handler interface.
type Handler interface {
	Handle(ev *Message) error
}
