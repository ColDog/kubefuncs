package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/coldog/kubefuncs/clients/go"
)

var handlers = map[string]kubefuncs.Handler{
	"responder": responder{},
	"archiver":  archiver{},
}

type archiver struct{}

func (h archiver) Handle(ev *kubefuncs.Message) error {
	log.Printf("%+v", ev)
	return nil
}

type responder struct{}

func (h responder) Handle(ev *kubefuncs.Message) error {
	return ev.Respond(&kubefuncs.HTTPResponse{
		Body: []byte("pong\n"),
	})
}

func main() {
	var handler string
	flag.StringVar(&handler, "handler", "responder", "handler to run (responder, archiver)")
	flag.Parse()

	err := kubefuncs.Run(handlers[handler])
	if err != nil {
		fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
		os.Exit(1)
	}
}
