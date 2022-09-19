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
