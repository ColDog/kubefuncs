package kubefuncs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	c, err := NewClient(WithCallEnabled())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer c.Close()

	c.On("test", "default", HandlerFunc(func(m *Message) error {
		return m.Respond(&Empty{})
	}))

	c.On("error", "default", HandlerFunc(func(m *Message) error {
		return errors.New("error")
	}))

	t.Run("Healthz", func(t *testing.T) {
		c.ListenHealthz()
		time.Sleep(200 * time.Millisecond)
		r, err := http.Get("http://127.0.0.1:8080/healthz")
		require.NoError(t, err)
		require.Equal(t, 200, r.StatusCode)
	})

	t.Run("Call", func(t *testing.T) {
		e, err := NewEvent("test", &Empty{})
		require.NoError(t, err)

		res, err := c.Call(ctx, e)

		require.NoError(t, err)
		require.Equal(t, res.Id, e.Id)
	})

	t.Run("CallError", func(t *testing.T) {
		e, err := NewEvent("error", &Empty{})
		require.NoError(t, err)

		res, err := c.Call(ctx, e)

		require.NoError(t, err)
		require.Equal(t, res.Id, e.Id)
		fmt.Printf(">> %+v\n", res)
	})
}
