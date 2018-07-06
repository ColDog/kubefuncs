# Local registry should be proxying into minikube.
REGISTRY := "localhost:5000"

# Allocate a random tag to always rebuild.
TAG := tag$(shell date +%s)

deploy-example:
	docker build -t localhost:5000/example:$(TAG) -f gateway/Dockerfile .
	docker push localhost:5000/example:$(TAG)
	@helm template charts/gateway --set=image.tag=$(TAG) --release kubefuncs --namespace kubefuncs > deploy.yaml
	kubectl apply -f deploy.yaml
	@rm deploy.yaml

deploy-nsq: render-charts
	kubectl apply -f charts/nsq/rendered.yaml

deploy-gateway:
	docker build -t localhost:5000/gateway:$(TAG) -f gateway/Dockerfile .
	docker push localhost:5000/gateway:$(TAG)
	@helm template charts/gateway --set=image.tag=$(TAG) --release kubefuncs --namespace kubefuncs > deploy.yaml
	kubectl apply -f deploy.yaml
	@rm deploy.yaml

build:
	@mkdir -p bin
	@go build -o bin/gateway ./gateway
	@go build -o bin/http-pong ./example/http-pong

release:
	docker build -t coldog/kubefuncs-gateway:$(GATEWAY_VERSION) -f gateway/Dockerfile .
	docker push coldog/kubefuncs-gateway:$(GATEWAY_VERSION)
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --all
