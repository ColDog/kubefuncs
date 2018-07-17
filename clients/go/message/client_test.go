package message

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testPubSub(t *testing.T, c Client) {
	c.Publish("test", []byte("hello"))

	msg := make(chan []byte, 1)
	c.Subscribe("test", "default", func(m []byte) error {
		msg <- m
		return nil
	})

	c.Publish("test", []byte("hello"))

	m := <-msg
	require.Equal(t, "hello", string(m))
}

func TestMockClient(t *testing.T) {
	testPubSub(t, NewMockClient())
}

func TestNSQClient(t *testing.T) {
	c, err := NewNSQClient("127.0.0.1:4150", "127.0.0.1:4161", nil)
	require.NoError(t, err)
	testPubSub(t, c)
}
