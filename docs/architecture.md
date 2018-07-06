# Architecture

## Functions

Functions are programs that listen to NSQ topics, do work, potentially respond or return errors. Invoking a function involves pushing an event into NSQ for a given topic.

If a return parameter is present in the message, the function should push a response onto the provided queue.

Functions are built using client libraries that support this minimal workflow and a set of common configuration. This allows for instrumentation in many different languages. A goal of this project is to avoid custom runtimes or other setup that moves away from core Kubernetes concepts.

## Deployment

Any Kubernetes resource can be used for a function. As long as a container is running in the cluster and able to connect to NSQ, we can talk to the function.

Recommended, and the default configuration supports this, is using a Deployment resource with autoscaling setup.

`kubefunctl` handles creating an opinionated deployment using a yaml configuration file for your application. This is a great way to get started but it's not required. As long as a function is listening to an NSQ topic, events can be routed to it.

## Metrics

Metrics are exported in Prometheus format by client libraries and NSQ.

## Gateway

Listens to HTTP requests and dispatches a request to a specific topic. This is a default component but is built on top of the core client libraries.
