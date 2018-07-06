# Getting Started

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
