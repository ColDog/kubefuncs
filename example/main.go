package main

import (
	"fmt"
	"os"

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
