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
		configFile string
	)
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "gateway listen address")
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
