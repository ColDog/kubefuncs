package kubefuncs

import (
	"log"
	"os"

	nsq "github.com/nsqio/go-nsq"
)

type nsqClient interface {
	publish(topic string, data []byte) error
	subscribe(topic, channel string, handler nsq.Handler) error
	close()
}

func newNsqClient(nsqdURL, lookupdURL string) (nsqClient, error) {
	c := &defaultNsqClient{nsqdURL: nsqdURL, lookupdURL: lookupdURL}
	if err := c.setupProducer(); err != nil {
		return nil, err
	}
	return c, nil
}

type defaultNsqClient struct {
	nsqdURL    string
	lookupdURL string
	producer   *nsq.Producer
	consumers  []*nsq.Consumer
}

func (h *defaultNsqClient) setupProducer() error {
	producer, err := nsq.NewProducer(h.nsqdURL, nsq.NewConfig())
	if err != nil {
		return err
	}
	producer.SetLogger(log.New(os.Stderr, "[nsq] ", log.LstdFlags), nsq.LogLevelDebug)
	h.producer = producer
	return nil
}

func (h *defaultNsqClient) publish(topic string, data []byte) error {
	return h.producer.PublishAsync(topic, data, nil)
}

func (h *defaultNsqClient) subscribe(topic, channel string, handler nsq.Handler) error {
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

func (h *defaultNsqClient) close() {
	log.Println("closing client")
	if h.producer != nil {
		h.producer.Stop()
	}
	for _, c := range h.consumers {
		c.Stop()
		<-c.StopChan
	}
}
