# Local registry should be proxying into minikube.
REGISTRY := "localhost:5000"

# Allocate a random tag to always rebuild.
TAG := t$(shell date +%s)

# Local setup:
# ===========

local/deploy-example:
	docker build -t localhost:5000/example:$(TAG) -f example/Dockerfile .
	docker push localhost:5000/example:$(TAG)
	helm upgrade --install \
		--values=example/functions.yaml \
		--set="image.repository=localhost:5000/example" \
		--set="image.tag=$(TAG)" \
		--namespace default \
		example charts/function

local/deploy-kubefuncs:
	docker build -t localhost:5000/gateway:$(TAG) -f gateway/Dockerfile .
	docker push localhost:5000/gateway:$(TAG)
	helm dep update charts/kubefuncs
	helm upgrade --install \
		--set="gateway.image.repository=localhost:5000/gateway" \
		--set="gateway.image.tag=$(TAG)" \
		--namespace kubefuncs \
		kubefuncs charts/kubefuncs
