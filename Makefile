# Local registry should be proxying into minikube.
REGISTRY := "localhost:5000"

# Allocate a random tag to always rebuild.
TAG := tag$(shell date +%s)

GATEWAY_VERSION := $(shell cat charts/kubefuncs/dependencies/gateway/Chart.yaml | grep 'version' | awk '{print $$2}')
VERSION := $(shell cat charts/kubefuncs/Chart.yaml | grep 'version' | awk '{print $$2}')

deploy-example:
	@./kubefunctl --build --apply \
		--image localhost:5000/example-app:$(TAG) \
		--namespace default \
		--release example-app \
		--functions example/functions.yaml \
		--build-dockerfile example/Dockerfile

deploy-nsq: render-charts
	kubectl apply -f charts/nsq/rendered.yaml

deploy-gateway:
	docker build -t localhost:5000/gateway:$(TAG) -f gateway/Dockerfile .
	docker push localhost:5000/gateway:$(TAG)
	@helm-template charts/gateway --set=tag=$(TAG) --release kubefuncs --namespace kubefuncs > deploy.yaml
	kubectl apply -f deploy.yaml
	@rm deploy.yaml

render-charts:
	helm-template charts/nsq --release kubefuncs --namespace kubefuncs > charts/nsq/rendered.yaml
	helm-template charts/gateway --release kubefuncs --namespace kubefuncs > charts/gateway/rendered.yaml

build:
	@mkdir -p bin
	@go build -o bin/gateway ./gateway
	@go build -o bin/http-pong ./example/http-pong

release:
	docker build -t coldog/kubefuncs-gateway:$(GATEWAY_VERSION) -f gateway/Dockerfile .
	docker push coldog/kubefuncs-gateway:$(GATEWAY_VERSION)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --all
