# KubeFuncs

Simple function architecture for Kubernetes.

Goals:
- Provide a low touch and extensible function framework.
- Allow for disparate and custom application architectures.
- Built on top of Docker and core Kubernetes resources.
- Simple opinionated structure.

This project came out of the desire to have a simple programming model for Kubernetes. Writing helm charts and deploying that through a CI system and wiring up services is incredibly powerful but very confusing for a newcomer, Kubernetes concepts are hard. This is an attempt at an opinionated application framework that supports multiple languages and is based on core Kubernetes concepts. The main actor is a control loop the listens to a queue and responds to events.

## Getting Started

Prerequisite: A working Kubernetes installation like minikube.

1. Clone the repo and setup core components.

```sh
git clone https://github.com/ColDog/kubefuncs
cd kubefuncs

kubectl apply -f charts/nsq/rendered.yaml
kubectl apply -f charts/gateway/rendered.yaml
```

2. Using the CLI, apply the example app into the cluster.

```sh
./kubefunctl --apply \
  --image TODO \
  --namespace default \
  --release example-app \
  --functions example/functions.yaml
```

3. Proxy the gateway locally

```sh
GATEWAY_POD=$(kubectl -n kubefuncs get pods | grep gateway | awk '{print $1}')
kubectl -n kubefuncs port-forward $GATEWAY_POD 8080:8080
```

4. Curl the gateway and receive a pong from the example app.
```sh
curl localhost:8080/test/hello
> pong
```

## Architecture

### Functions

Functions are programs that listen to NSQ topics, do work, potentially respond or return errors. Invoking a function involves pushing an event into NSQ for a given topic.

If a return parameter is present in the message, the function should push a response onto the provided queue.

### Deployment

Any Kubernetes resource can be used for a function. As long as a container is running in the cluster and able to connect to NSQ, we can talk to the function.

Recommended, and the default configuration supports this, is using a Deployment resource with autoscaling setup.

`kubefunctl` handles creating an opinionated deployment using a Procfile for your application. This is a great way to get started but it's not required.

Kubernetes is also not a dependency, despite the name, significantly more infrastructure will be required to build KubeFunc as it relies heavily on Kubernetes resources.

### Metrics

Metrics are exported in Prometheus format by client libraries and NSQ.

### Components

- Gateway: Listens to HTTP requests and dispatches a request to a specific topic. This is a default component but is built on top of the core client libraries.
