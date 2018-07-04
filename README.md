# KubeFuncs

Simple function architecture for Kubernetes.

Goals:
- Provide a low touch and extensible function framework.
- Allow for disparate and custom application architectures.
- Built on top of Docker and core Kubernetes resources.
- Simple opinionated structure.

## Getting Started

We require Helm to deploy resources into Kubernetes. You can also apply using helm-template if Helm cannot be used.

TODO setup instructions, separate from this repository.

Custom application resources can be generated as well, KubeFuncs uses core Kubernetes resources in v1 status, use the generator to get started but customize as needed.

Projects are generated for specific languages using client libraries.

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
