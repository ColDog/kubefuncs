package message

// Handler represents a message handler.
type Handler func(msg []byte) error

// Client represents a messaging client.
type Client interface {
	Publish(topic string, data []byte) error
	Subscribe(topic, channel string, handler Handler) error
	Close()
}

// NewMockClient returns an in memory only client.
func NewMockClient() Client { return &mockClient{handlers: map[string][]Handler{}} }

type mockClient struct {
	handlers map[string][]Handler
}

func (m *mockClient) Publish(topic string, data []byte) error {
	for _, handler := range m.handlers[topic] {
		go handler(data)
	}
	return nil
}

func (m *mockClient) Subscribe(topic, channel string, handler Handler) error {
	m.handlers[topic] = append(m.handlers[topic], handler)
	return nil
}

func (m *mockClient) Close() {}
