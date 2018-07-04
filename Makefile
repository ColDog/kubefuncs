REGISTRY := "localhost:5000"
TAG := tag$(shell date +%s)

deploy-example:
	@./deploy --build --apply \
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
