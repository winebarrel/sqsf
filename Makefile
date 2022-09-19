.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/sqsf

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf sqsf sqsf.exe
