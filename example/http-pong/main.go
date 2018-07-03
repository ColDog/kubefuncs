package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coldog/kubefuncs/clients/go"
)

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
	os.Exit(1)
}

type handler struct{}

func (h handler) Handle(ev *kubefuncs.Message) error {
	resp, err := kubefuncs.NewEvent(ev.Event.Return, &kubefuncs.HTTPResponse{
		Body: []byte("pong"),
	})
	if err != nil {
		return err
	}
	ev.Respond(resp)
	return nil
}

func main() {
	var (
		lookupdURL string
		nsqdURL    string
		topic      string
	)
	flag.StringVar(&topic, "topic", "test", "topic to listen on")
	flag.StringVar(&lookupdURL, "lookupd-addr", "127.0.0.1:4161", "nsqlookupd address")
	flag.StringVar(&nsqdURL, "nsqd-addr", "127.0.0.1:4150", "nsqd address")
	flag.Parse()

	client, err := kubefuncs.NewClient(
		kubefuncs.WithLookupdURL(lookupdURL),
		kubefuncs.WithNsqdURL(nsqdURL),
	)
	if err != nil {
		handleErr(err)
	}

	client.On(topic, "default", handler{})
	client.Wait()
}
