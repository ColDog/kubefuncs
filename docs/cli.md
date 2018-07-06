# CLI

This is the CLI that supports building and deploying functions into a Kubernetes cluster. It depends on the docker daemon for building and kubectl for deployment.

```sh
$ ./kubefunctl -h
Usage: kubefunctl [options]
        --image IMAGE                docker image url
        --namespace NAMESPACE        kubernetes namespace
        --release RELEASE            release name
        --build-context CONTEXT      docker build context
        --build-dockerfile FILE      dockerfile path
        --build                      build the docker container
        --functions FILE             functions.yaml file
        --apply                      apply manifests using kubectl
```

For example, to build and push the example application to a local docker registry and deploy it into the cluster, we need to run the following:

```sh
kubefunctl --build --apply \
  --image localhost:5000/example-app:$(TAG) \
  --namespace default \
  --release example-app \
  --functions example/functions.yaml \
  --build-dockerfile example/Dockerfile
```

To deploy, the cli will build a directory `manifests/` with the necessary manifests needed to deploy to the cluster and then run `kubectl apply -f manifests`. Functions are configured via a configuration file:

```yaml
test: # Function Name.
  command: ['/bin/http-pong'] # Command to execute the function.
  env: # Kubernetes environment variables.
  - name: test
    value: test
  ports: # Kubernetes ports.
  - name: http
    containerPort: 8080
```

All functions result in Kubernetes deployment resources and can be autoscaled using the built in horizontal pod autoscaler.
