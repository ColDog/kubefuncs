package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/coldog/kubefuncs/clients/go"
)

func handleErr(err error) {
	fmt.Fprintf(os.Stderr, "exit 1: %v\n", err)
	os.Exit(1)
}

func main() {
	var (
		listenAddr string
		lookupdURL string
		nsqdURL    string
		gatewayID  string
		configFile string
	)
	hostname, _ := os.Hostname()
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "gateway listen address")
	flag.StringVar(&gatewayID, "gateway-id", hostname, "gateway id is a unique id for this gateway instance")
	flag.StringVar(&lookupdURL, "lookupd-addr", "127.0.0.1:4161", "nsqlookupd address")
	flag.StringVar(&nsqdURL, "nsqd-addr", "127.0.0.1:4150", "nsqd address")
	flag.StringVar(&configFile, "config", "routes.json", "routes configuration file")
	flag.Parse()

	routes := map[string]string{}
	f, err := os.Open(configFile)
	if err != nil {
		handleErr(err)
	}
	err = json.NewDecoder(f).Decode(&routes)
	if err != nil {
		handleErr(err)
	}

	client, err := kubefuncs.NewClient(
		kubefuncs.WithCallEnabled(),
		kubefuncs.WithClientID(gatewayID),
		kubefuncs.WithLookupdURL(lookupdURL),
		kubefuncs.WithNsqdURL(nsqdURL),
	)
	if err != nil {
		handleErr(err)
	}

	router := &Router{
		Client: client,
		Routes: routes,
	}

	err = http.ListenAndServe(listenAddr, router)
	if err != nil {
		handleErr(err)
	}
}
