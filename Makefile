NAME := prometheus-kafka-adapter
PACKAGE_NAME := github.com/Telefonica/prometheus-kafka-adapter
GO_VER := 1.17.5
LIBC_GO_VER := $(GO_VER)-buster
MUSL_GO_VER := $(GO_VER)-alpine

all: fmt test build

fmt:
	docker run --rm -v $(CURDIR):/app:z -w /app golang:$(MUSL_GO_VER) gofmt -l -w -s *.go

test:
	docker run --rm -v $(CURDIR):/app:z -w /app golang:$(MUSL_GO_VER) sh tools/testscript.sh vet
	docker run --rm -v $(CURDIR):/app:z -w /app golang:$(MUSL_GO_VER) sh tools/testscript.sh test

build-docker-image:
	docker build -t telefonica/prometheus-kafka-adapter:latest .
