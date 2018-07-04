package kubefuncs

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"
)

type Message struct {
	*Event
	Response *Event
}

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

func (m *Message) Payload(dst proto.Message) error {
	return ptypes.UnmarshalAny(m.Event.Payload, dst)
}

func (m *Message) Respond(payload proto.Message) error {
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

type Handler interface {
	Handle(ev *Message) error
}
