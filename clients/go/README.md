# Client Go

A client library for Kubefuncs in go.

The simplest example is:

```go
// main.go
package main

import (
  "fmt"

  client "github.com/coldog/kubefuncs/clients/go"
)

func main() {
  err := client.Run(client.HandlerFunc(func(ev *client.Message) error {
    return ev.Respond(&client.HTTPResponse{
      Body: []byte("pong\n"),
    })
  }))
  if err != nil {
    fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
    os.Exit(1)
  }
}
```

For more advanced architectures, instantiate a client:

```go
// Initiate a new client.
c, err := client.NewClient()

// A new event must be initialized by calling NewEvent(...).
e, err := client.NewEvent("test", &client.HTTPRequest{})

// Emit will publish without waiting for a response.
c.Emit("test", e)

// Call will publish and wait for a response.
c.Call("test", e)

// On will register a handler.
c.On("test", "default", client.HandlerFunc(func (m *client.Message) error {
  return m.Respond(&client.Empty{})
}))
```
