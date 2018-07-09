# example

This is an example application that receives HTTP requests and sends responses.

Read the getting started [guide](../#getting-started) for more information.

## Deployment

Deploy the functions.

```bash
helm upgrade --install \
  --values functions.yaml \
  --namespace default \
  example "kubefuncs/function"
```

Deploy the gateway.

```bash
helm upgrade --install \
  --namespace kubefuncs \
  --values gateway.yaml \
  kubefuncs kubefuncs/kubefuncs
```
