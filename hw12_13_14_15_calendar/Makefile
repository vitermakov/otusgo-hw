BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	docker-compose -f ./deployments/rundb.yaml build
	docker-compose -f ./deployments/docker-compose.yaml build

run: build
	docker-compose -f ./deployments/rundb.yaml up -d
	docker-compose -f ./deployments/docker-compose.yaml up -d

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	rm -Rf internal/handler/grpc/pb
	mkdir internal/handler/grpc/pb
	protoc --proto_path=internal/handler/grpc/proto/ --go_out=. --go-grpc_out=. internal/handler/grpc/proto/*.proto

.PHONY: build run build-img run-img version migrate test lint