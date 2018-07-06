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

// Option represents a configuration parameter.
type Option func(o *opts)

// WithLookupdURL configures the lookupd instance url. This defaults to the env
// variable $NSQ_LOOKUPD_ADDR:$NSQ_LOOKUPD_PORT.
func WithLookupdURL(url string) Option { return func(o *opts) { o.lookupdURL = url } }

// WithNsqdURL configures the lookupd instance url. This defaults to the env
// variable $NSQ_NSQD_ADDR:$NSQ_NSQD_PORT.
func WithNsqdURL(url string) Option { return func(o *opts) { o.nsqdURL = url } }

// WithClientID configures the unique client id for this instance, it defaults
// to the hostname.
func WithClientID(id string) Option { return func(o *opts) { o.clientID = id } }

// WithCallEnabled ensures that this client can handle responses from published
// events. This must be enabled to use the Call(...) method.
func WithCallEnabled() Option { return func(o *opts) { o.rpc = true } }

// WithTopic configures the new topic.
func WithTopic(topic string) Option { return func(o *opts) { o.topic = topic } }

// WithChannel configures the channel.
func WithChannel(channel string) Option { return func(o *opts) { o.channel = channel } }

// WithHealthzAddr adds the default healthz address.
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
