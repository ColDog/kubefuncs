package message

import (
	"log"

	nsq "github.com/nsqio/go-nsq"
)

// NewNSQClient returns a new Client interface using NSQ as the transport.
func NewNSQClient(nsqdURL, lookupdURL string, logger *log.Logger) (Client, error) {
	c := &nsqClient{nsqdURL: nsqdURL, lookupdURL: lookupdURL}
	if err := c.setupProducer(); err != nil {
		return nil, err
	}
	return c, nil
}

type nsqClient struct {
	nsqdURL    string
	lookupdURL string
	logger     *log.Logger
	producer   *nsq.Producer
	consumers  []*nsq.Consumer
}

func (h *nsqClient) setupProducer() error {
	producer, err := nsq.NewProducer(h.nsqdURL, nsq.NewConfig())
	if err != nil {
		return err
	}
	if h.logger != nil {
		producer.SetLogger(h.logger, nsq.LogLevelInfo)
	}
	h.producer = producer
	return nil
}

func (h *nsqClient) Publish(topic string, data []byte) error {
	return h.producer.PublishAsync(topic, data, nil)
}

func (h *nsqClient) Subscribe(topic, channel string, handler Handler) error {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}
	if h.logger != nil {
		consumer.SetLogger(h.logger, nsq.LogLevelDebug)
	}
	consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		return handler(m.Body)
	}))
	h.consumers = append(h.consumers, consumer)
	err = consumer.ConnectToNSQLookupd(h.lookupdURL)
	if err != nil {
		return err
	}
	return nil
}

func (h *nsqClient) Close() {
	log.Println("closing client")
	if h.producer != nil {
		h.producer.Stop()
	}
	for _, c := range h.consumers {
		c.Stop()
		<-c.StopChan
	}
}
