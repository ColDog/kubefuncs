# Gateway

The gateway package translates http calls to function calls and then takes the response and returns them to the client. There are standard facilities for processing HTTP calls in all client libraries.

Read through the getting started guide to see an example of setting up the gateway.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

To install the chart with the release name `my-release`, run:

```bash
$ helm repo add kubefuncs https://s3.amazonaws.com/kubefuncs-chart-repository
$ helm install --name my-release kubefuncs/gateway
```

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The core configuration parameters needed for the gateway is the route configuration. The example below will route any calls to `/test/*` to the `test` function:

```yaml
config:
  routes:
    /test/: test
```

The routes key contains path prefixes to the path that they should be routed to.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
$ helm install --name my-release \
    --set rbac.create=true \
    kubefuncs/gateway
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml kubefuncs/gateway
```

> **Tip**: You can use the default [values.yaml](values.yaml)
