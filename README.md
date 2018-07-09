# Kubefuncs

Building blocks for lightweight functions on top of Kubernetes. No custom runtimes, focuses on using docker and core Kubernetes resources.

The kubefuncs project is a set of helm charts and client libraries that together provide a toolkit for developing functions. A function is a program running a client library that listens for events along a queue. The goal is to keep client libraries and the runtime simple so that we can use core kubernetes resources.

This project is currently in **alpha** status.

## Contents

* [Getting Started](#getting-started)
* [Architecture](#architecture)
* [Charts](charts)
* [Clients](clients)
* [Example](example)

## Getting Started

Requirements:
- A running Kubernetes cluster.
- Helm cli installed.

First step is to get the core Kubefuncs dependencies installed into your cluster, this consists of an NSQ installation and the Kubefuncs gateway to handle http events.

To do this, add the helm repository:

```bash
CHART_REPO="https://s3.amazonaws.com/kubefuncs-chart-registry"
helm repo add kubefuncs ${CHART_REPO}
```

Now, we should be able to install the main Kubefuncs package which contains both NSQ and the gateway:

```bash
helm upgrade --install --namespace kubefuncs kubefuncs kubefuncs/kubefuncs
```

You can check on resources in the kubefuncs namespace:

```bash
kubectl -n kubefuncs get pods
```

Now for the application setup. First step is to setup an application with a dockerfile. The application should use the kubefuncs client library for the specific language. An example dockerfile and main function for a go application can be setup with the following:

```dockerfile
# Dockerfile
FROM golang:1.10-alpine as builder
COPY . /go/src/github.com/coldog/kubefuncs/example
RUN go get ./... \
  && go build \
    -o /build/app \
    github.com/coldog/kubefuncs/example

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/app /bin/app
CMD /bin/app
```

And the main function we only need:

```golang
// main.go
package main

import (
  "fmt"

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
```

We can now build the project just using `docker build -t <my-tag> .`. If you have a local registry in your Kubernetes cluster to push to that will work, if not you can push to docker hub or somewhere else.

Once we have the image available to our Kubernetes cluster we can deploy. To deploy, setup a `functions.yaml` file which describes which functions we want to run:

```yaml
# functions.yaml
image:
  repository: <my-repo>
  tag: <my-tag>

functions:
- name: test
```

This sets up a single function named `test`. We can now deploy using helm:

```bash
helm upgrade --install \
  --values functions.yaml \
  --namespace default \
  example "kubefuncs/function"
```

This command installs a new helm package using the `kubefuncs/function` chart. We use the `functions.yaml` file as the config value for the chart. The chart will loop through all of the functions and create the appropriate Kubernetes resources.

Let's now test the gateway, since our function is setup to return `pong` if we curl the gateway we should be able to receive a pong response. First, port forward the gateway to your local machine.

```bash
POD=$(kubectl -n kubefuncs get pods | grep gateway | awk '{print $1}')
kubectl port-forward $POD 8080
```

Now, in a new terminal, we should be able to curl the gateway and receive the pong response:

```bash
curl localhost:8080/test/hello
> pong
```

By default, the gateway is configured to send the `/test/*` path to the test function. This can be configured by updating the helm chart gateway. Let's create a new config file to add a new route:

```yaml
# gateway.yaml
gateway:
  config:
    routes:
      /ping/: test
```

This config file maps path prefixes to topics which call functions. When we added our `test` function it was configured to listen on the test topic. The routes file that was added will now send the `/ping/` path along to the test function. More information on this is available in the [helm chart](charts/gateway).

Now, we can update our gateway:

```bash
helm upgrade --install \
  --namespace kubefuncs \
  --values gateway.yaml \
  kubefuncs kubefuncs/kubefuncs
```

And then let's test the new route:

```bash
curl localhost:8080/ping/hello
> pong
```

We've now completed the core concepts around Kubefuncs. You should now be able to author more complex packages and compose more complex architectures.

## Architecture

### Functions

Functions are programs that listen to NSQ topics, do work, potentially respond or return errors. Invoking a function involves pushing an event into NSQ for a given topic.

If a return parameter is present in the message, the function should push a response onto the provided queue.

Functions are built using client libraries that support this minimal workflow and a set of common configuration. This allows for instrumentation in many different languages. A goal of this project is to avoid custom runtimes or other setup that moves away from core Kubernetes concepts.

### Deployment

Any Kubernetes resource can be used for a function. As long as a container is running in the cluster and able to connect to NSQ, we can talk to the function.

The [function](charts/function) chart sets up the default configuration. This is recommended but not required.

### Metrics

Metrics are exported in Prometheus format by client libraries and NSQ.

### Gateway

Listens to HTTP requests and dispatches a request to a specific topic. This is a default component but is built on top of the core client libraries.

## Developing

Required tools:

* `make`
* `minikube`
* `wrk`
* `docker`

Run `make local/setup` to get your minikube installation going. You can then run `make local/deploy-kubefuncs` and `make local/deploy-example` to deploy the kubefuncs package and example app. Once this is setup, `make test/e2e` runs some basic e2e tests.
