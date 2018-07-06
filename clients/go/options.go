package kubefuncs

import "os"

func env(name, defaults string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaults
	}
	return val
}

func defaults() *opts {
	hostname, _ := os.Hostname()
	return &opts{
		topic:   env("TOPIC", "test"),
		channel: env("CHANNEL", "default"),
		lookupdURL: env("NSQ_LOOKUPD_ADDR", "127.0.0.1") +
			":" + env("NSQ_LOOKUPD_PORT", "4161"),
		nsqdURL: env("NSQ_NSQD_ADDR", "127.0.0.1") +
			":" + env("NSQ_NSQD_PORT", "4150"),
		clientID:    hostname,
		healthzAddr: env("HEALTHZ_ADDR", ":8080"),
	}
}

type Option func(o *opts)

func WithLookupdURL(url string) Option   { return func(o *opts) { o.lookupdURL = url } }
func WithNsqdURL(url string) Option      { return func(o *opts) { o.nsqdURL = url } }
func WithClientID(id string) Option      { return func(o *opts) { o.clientID = id } }
func WithCallEnabled() Option            { return func(o *opts) { o.rpc = true } }
func WithTopic(topic string) Option      { return func(o *opts) { o.topic = topic } }
func WithChannel(channel string) Option  { return func(o *opts) { o.channel = channel } }
func WithHealthzAddr(addr string) Option { return func(o *opts) { o.healthzAddr = addr } }

type opts struct {
	lookupdURL  string
	nsqdURL     string
	clientID    string
	topic       string
	channel     string
	rpc         bool
	healthzAddr string
}
