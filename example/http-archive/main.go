package main

import (
	"fmt"
	"log"
	"os"

	"github.com/coldog/kubefuncs/clients/go"
)

type handler struct{}

func (h handler) Handle(ev *kubefuncs.Message) error {
	log.Printf("%+v", ev)
	return nil
}

func main() {
	// Different channel than the main pong application.
	err := kubefuncs.Run(handler{}, kubefuncs.WithChannel("archive"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
		os.Exit(1)
	}
}
