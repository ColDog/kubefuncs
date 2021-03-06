# Kubefuncs

Building blocks for lightweight functions on top of Kubernetes. No custom runtimes, focuses on using docker and core Kubernetes resources.

The kubefuncs project is a set of helm charts and client libraries that together provide a toolkit for developing functions. A function is a program running a client library that listens for events along a queue. The goal is to keep client libraries and the runtime simple so that we can use core kubernetes resources.

This project is currently in **alpha** status.

## Contents

* [Getting Started](example/README.md)
* [Architecture](#architecture)
* [Charts](charts)
* [Clients](clients)

## Overview

Functions are programs that listen to NSQ topics, do work, potentially respond or return errors. Invoking a function involves pushing an event into NSQ for a given topic.

If a return parameter is present in the message, the function should push a response onto the provided queue.

Functions are built using client libraries that support this minimal workflow and a set of common configuration. This allows for instrumentation in many different languages. A goal of this project is to avoid custom runtimes or other setup that moves away from core Kubernetes concepts.

Any Kubernetes resource can be used for a function. As long as a container is running in the cluster and able to connect to NSQ, we can talk to the function.

The [function](charts/function/README.md) chart sets up the default configuration. This is recommended but not required.

## Installation

Requires:

- A working Kubernetes cluster.
- A helm installation.

Add the helm repository:

```bash
CHART_REPO="https://charts.kubefuncs.com"
helm repo add kubefuncs ${CHART_REPO}
```

Deploy the gateway and NSQ.

```bash
helm upgrade --install \
  --namespace kubefuncs \
  --values gateway.yaml \
  kubefuncs kubefuncs/kubefuncs
```

## Developing

Required tools:

* `make`
* `minikube`
* `wrk`
* `docker`

Run `make local/setup` to get your minikube installation going. You can then run `make local/deploy-kubefuncs` and `make local/deploy-example` to deploy the kubefuncs package and example app. Once this is setup, `make test/e2e` runs some basic e2e tests.

### Releasing

1. Update the chart versions in `chart/<name>/Chart.yaml` that need to be updated.
2. Update the chart requirements in `chart/kubefuncs/requirements.yaml` that need to be updated.
3. Update `_coverpage.yaml` if the kubefuncs chart was updated.
4. Run `make test/e2e` to be sure release hasn't broken anything.
5. Run `make release`.
