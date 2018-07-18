package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	kubefuncs "github.com/coldog/kubefuncs/clients/go"
	"github.com/stretchr/testify/require"
)

func TestGateway(t *testing.T) {
	client, err := kubefuncs.NewClient(
		kubefuncs.WithCallEnabled(),
		kubefuncs.WithMockClient(),
	)
	require.NoError(t, err)

	router := &Router{
		Config: Config{
			Routes: map[string]string{
				"/test/": "test",
			},
		},
		Client: client,
	}

	go func() {
		err := http.ListenAndServe(":8081", router)
		require.Nil(t, err)
	}()
	time.Sleep(200 * time.Millisecond)

	client.On("test", "default", kubefuncs.HandlerFunc(func(ev *kubefuncs.Message) error {
		return ev.Respond(&kubefuncs.HTTPResponse{
			Body: []byte("pong\n"),
		})
	}))

	t.Run("Ping", func(t *testing.T) {
		r, err := http.Get("http://127.0.0.1:8081/test/ping")
		require.NoError(t, err)

		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		require.Equal(t, "pong\n", string(data))
	})
}
