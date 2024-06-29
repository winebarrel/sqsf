AWS_REGION              := us-east-1
AWS_ENDPOINT_URL        := http://localhost:4566
DOCKER_AWS_ENDPOINT_URL := http://localstack:4566
QUEUE_NAME       := sqsf-test
QUEUE_URL        := http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/sqsf-test

.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/sqsf

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: clean
clean:
	rm -rf sqsf sqsf.exe

.PHONY: queue
queue:
	docker compose run awscli sqs --region $(AWS_REGION) --endpoint-url $(DOCKER_AWS_ENDPOINT_URL) \
		create-queue --queue-name $(QUEUE_NAME)

.PHONY: message
message:
	docker compose run awscli sqs --region $(AWS_REGION) --endpoint-url $(DOCKER_AWS_ENDPOINT_URL) \
		send-message --queue-url $(QUEUE_URL) --message-body 'hello world'

.PHONE: receive
receive:
	go run ./cmd/sqsf --region $(AWS_REGION) --endpoint-url $(AWS_ENDPOINT_URL) $(QUEUE_NAME)
