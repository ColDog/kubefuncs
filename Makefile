GATEWAY_VERSION := $(shell cat charts/gateway/Chart.yaml | grep version | awk '{print $$2}')
FUNCTION_VERSION := $(shell cat charts/function/Chart.yaml | grep version | awk '{print $$2}')
KUBEFUNCS_VERSION := $(shell cat charts/kubefuncs/Chart.yaml | grep version | awk '{print $$2}')
NSQ_VERSION := $(shell cat charts/nsq/Chart.yaml | grep version | awk '{print $$2}')

CHART_BUCKET := kubefuncs-chart-repository
CHART_REPO := https://s3.amazonaws.com/$(CHART_BUCKET)

export AWS_PROFILE=kubefuncs
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
		--values=example/gateway.yaml \
		--set="gateway.image.repository=localhost:5000/gateway" \
		--set="gateway.image.tag=$(TAG)" \
		--namespace kubefuncs \
		kubefuncs charts/kubefuncs

local/setup:
	minikube start
	kubectl --context minikube apply -f tests/registry.yaml
	kubectl --context minikube apply -f tests/ingress.yaml
	kubectl label node minikube role=ingress-controller
	pod_id=$$(kubectl -n kube-system get pods | grep 'registry' | awk '{print $$1}') && \
		kubectl -n kube-system port-forward $$pod_id 5000:5000

define package
	helm init --client-only
	helm repo add kubefuncs $(CHART_REPO)

	mkdir -p .repo
	curl -o .repo/index.yaml $(CHART_REPO)/index.yaml || echo 'repo not initialized'

	echo "charts/$(1)"
	helm package -d .repo ./charts/$(1)
	helm repo index --url $(CHART_REPO) .repo

	aws s3 cp .repo/index.yaml s3://$(CHART_BUCKET)/index.yaml

	# TODO: Only copy if does not exist.
	aws s3 cp \
		.repo/$(1)-$(shell cat charts/$(1)/Chart.yaml | grep version | awk '{print $$2}').tgz \
		s3://$(CHART_BUCKET)/$(1)-$(shell cat charts/$(1)/Chart.yaml | grep version | awk '{print $$2}').tgz
endef

build/proto:
	cd clients/proto && ./build.sh

release/function:
	$(call package,function)

release/gateway:
	docker build -t coldog/kubefuncs-gateway:$(GATEWAY_VERSION) -f gateway/Dockerfile .
	docker tag coldog/kubefuncs-gateway:$(GATEWAY_VERSION) coldog/kubefuncs-gateway:latest
	docker push coldog/kubefuncs-gateway:$(GATEWAY_VERSION)
	docker push coldog/kubefuncs-gateway:latest
	$(call package,gateway)

release/example:
	docker build -t coldog/kubefuncs-example:latest -f example/Dockerfile .
	docker push coldog/kubefuncs-example:latest

release/nsq:
	$(call package,nsq)

release/kubefuncs:
	helm dep update charts/kubefuncs
	helm template charts/kubefuncs > charts/kubefuncs/bundle.yaml
	$(call package,kubefuncs)

release/git:
	git add -A
	git commit -m 'Release $(KUBEFUNCS_VERSION)'
	git tag -a $(KUBEFUNCS_VERSION) -m "Release $(KUBEFUNCS_VERSION)"
	git push --tags origin master

release/docs:
	aws s3 sync ./ s3://kubefuncs.com \
		--exclude "*" \
		--include "charts/*.md" \
		--include "clients/*.md" \
		--include "gateway/*.md" \
		--include "example/*.md" \
		--include "_coverpage.md" \
		--include "_sidebar.md" \
		--include "README.md" \
		--include "index.html"

release: release/function release/nsq release/gateway release/example release/kubefuncs release/docs release/git

test/e2e:
	@tests/e2e.sh
