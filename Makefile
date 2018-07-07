GATEWAY_VERSION := $(shell cat charts/gateway/Chart.yaml | grep version | awk '{print $$2}')
FUNCTION_VERSION := $(shell cat charts/function/Chart.yaml | grep version | awk '{print $$2}')
KUBEFUNCS_VERSION := $(shell cat charts/kubefuncs/Chart.yaml | grep version | awk '{print $$2}')
NSQ_VERSION := $(shell cat charts/nsq/Chart.yaml | grep version | awk '{print $$2}')

export AWS_PROFILE=personal
export AWS_REGION=us-east-1

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

release/gateway-docker:
	docker build -t coldog/kubefuncs-gateway:$(GATEWAY_VERSION) -f gateway/Dockerfile .
	docker tag coldog/kubefuncs-gateway:$(GATEWAY_VERSION) coldog/kubefuncs-gateway:latest
	docker push coldog/kubefuncs-gateway:$(GATEWAY_VERSION)
	docker push coldog/kubefuncs-gateway:latest

release/helm-init:
	helm s3 init s3://kubefuncs-chart-registry
	helm repo add kubefuncs s3://kubefuncs-chart-registry

release/function:
	helm package ./charts/function
	helm s3 push ./function-$(FUNCTION_VERSION).tgz kubefuncs

release/gateway: release/docker-gateway
	helm package ./charts/gateway
	helm s3 push ./gateway-$(GATEWAY_VERSION).tgz kubefuncs

release/example:
	docker build -t coldog/kubefuncs-example:latest -f example/Dockerfile .
	docker push coldog/kubefuncs-example:latest

release/nsq:
	helm package ./charts/nsq
	helm s3 push ./nsq-$(NSQ_VERSION).tgz nsq

release/kubefuncs-bundle:
	helm template charts/kubefuncs > charts/kubefuncs/bundle.yaml

release/kubefuncs: release/kubefuncs-bundle
	helm dep update ./charts/kubefuncs
	helm package ./charts/kubefuncs
	helm s3 push ./kubefuncs-$(KUBEFUNCS_VERSION).tgz nsq

	git commit -m "Release $(KUBEFUNCS_VERSION)"
	git tag -a $(KUBEFUNCS_VERSION) -m "Release $(KUBEFUNCS_VERSION)"
	git push --all
