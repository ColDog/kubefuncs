package main

import (
	"fmt"
	"os"

	"github.com/coldog/kubefuncs/clients/go"
)

type handler struct{}

func (h handler) Handle(ev *kubefuncs.Message) error {
	return ev.Respond(&kubefuncs.HTTPResponse{
		Body: []byte("pong\n"),
	})
}

func main() {
	err := kubefuncs.Run(handler{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
		os.Exit(1)
	}
}
