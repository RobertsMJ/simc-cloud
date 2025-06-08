IMAGE_NAME ?= echo-lambda
ECR_REPO ?= $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(IMAGE_NAME)
TAG ?= latest

.PHONY: build docker-build docker-tag docker-push sam-test

build:
	GOOS=linux GOARCH=amd64 go build -o main main.go

docker-build: build
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-tag:
	docker tag $(IMAGE_NAME):$(TAG) $(ECR_REPO):$(TAG)

docker-push: docker-build docker-tag
	docker push $(ECR_REPO):$(TAG)

sam-test: docker-build
	echo '{"body": "test input"}' > event.json
	sam local invoke --container-image $(IMAGE_NAME):$(TAG) -e event.json
	rm event.json
