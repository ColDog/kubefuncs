# kubefuncs

Installs all components needed to run Kubefuncs. This includes the following subcharts:

* [Gateway](../gateway)
* [NSQ](../nsq)

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

To install the chart with the release name `my-release`, run:

```bash
$ helm repo add kubefuncs https://s3.amazonaws.com/kubefuncs-chart-repository
$ helm install --name my-release kubefuncs/kubefuncs
```

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

View the [values.yaml](values.yaml) file for configuration information.

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
$ helm install --name my-release \
    --set rbac.create=true \
    kubefuncs/kubefuncs
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml kubefuncs/kubefuncs
```

> **Tip**: You can use the default [values.yaml](values.yaml)
