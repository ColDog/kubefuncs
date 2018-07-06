package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	kubefuncs "github.com/coldog/kubefuncs/clients/go"
)

// Config is the configuration for the router.
type Config struct {
	// Routes are a map of prefixes to topics.
	Routes map[string]string `json:"routes"`
}

// Router returns the given topic for the provided route.
type Router struct {
	Config
	Client *kubefuncs.Client
}

// Route will return the topic for the given path.
func (ro *Router) route(path string) string {
	if strings.HasPrefix(path, "/health") {
		return "health"
	}
	for key, val := range ro.Routes {
		if strings.HasPrefix(path, key) {
			return val
		}
	}
	return ""
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	log.Printf("handling: %s routes=%+v", r.URL.Path, ro.Routes)
	topic := ro.route(r.URL.Path)
	if topic == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if topic == "health" {
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ev, err := kubefuncs.NewEvent(topic, &kubefuncs.HTTPRequest{
		Url:  r.URL.String(),
		Body: body,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg, err := ro.Client.Call(ctx, ev)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	resp := &kubefuncs.HTTPResponse{}
	err = msg.Payload(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp.Body)
}
