package kubefuncs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEvent(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		_, err := NewEvent("", &Empty{})
		require.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		e, err := NewEvent("test", &Empty{})
		require.NoError(t, err)
		require.NotEmpty(t, e.Id)
		require.Empty(t, e.Return)
	})
}

func TestMessage(t *testing.T) {
	t.Run("Payload", func(t *testing.T) {
		e, _ := NewEvent("test", &Empty{})
		m := &Message{Event: e}

		o := &Empty{}
		err := m.Payload(o)
		require.NoError(t, err)
	})

	t.Run("RespondNoReturn", func(t *testing.T) {
		e, _ := NewEvent("test", &Empty{})
		m := &Message{Event: e}

		err := m.Respond(&Empty{})
		require.Error(t, err)
	})

	t.Run("Respond", func(t *testing.T) {
		e, _ := NewEvent("test", &Empty{})
		e.Return = "test2"
		m := &Message{Event: e}

		err := m.Respond(&Empty{})
		require.NoError(t, err)
	})
}
