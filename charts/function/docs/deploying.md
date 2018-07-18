# Deploying Functions


Deploying is done (as in the getting started guide) by using helm:

```bash
helm upgrade --install \
  --values functions.yaml \
  --namespace default \
  example "kubefuncs/function"
```

Note, if you don't use helm you can merely initialize helm using the client only flag:

```bash
helm init --client-only
```

And then just use `helm template` to render your functions and apply them:

```bash
helm template \
  --values functions.yaml \
  --namespace default \
  --release example \
  "kubefuncs/function" > my-yaml.yaml
kubectl apply -f my-yaml.yaml
```
