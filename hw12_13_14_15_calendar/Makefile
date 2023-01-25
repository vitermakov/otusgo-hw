BIN_CALENDAR=./bin/calendar
CFG_CALENDAR_FILE=./deployments/configs/calendar_config.json
CFG_CALENDAR_TPL=./deployments/configs/calendar_config.json.template
BIN_SCHEDULER=./bin/scheduler
CFG_SCHEDULER_FILE=./deployments/configs/scheduler_config.json
CFG_SCHEDULER_TPL=./deployments/configs/scheduler_config.json.template
BIN_SENDER=./bin/sender
CFG_SENDER_FILE=./deployments/configs/sender_config.json
CFG_SENDER_TPL=./deployments/configs/sender_config.json.template
ENV_NAME=tests

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

up: cfg-up
	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml up -d --build

down: cfg-clean
	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down

restart: down up

cfg-up:
	env `cat ./deployments/env.${ENV_NAME}` envsubst < ${CFG_CALENDAR_TPL} > ${CFG_CALENDAR_FILE}; \
	env `cat ./deployments/env.${ENV_NAME}` envsubst < ${CFG_SCHEDULER_TPL} > ${CFG_SCHEDULER_FILE}; \
	env `cat ./deployments/env.${ENV_NAME}` envsubst < ${CFG_SENDER_TPL} > ${CFG_SENDER_FILE};

cfg-clean:
	rm -f ${CFG_CALENDAR_FILE} ${CFG_SCHEDULER_FILE} ${CFG_SENDER_FILE}

integration-test: cfg-up
	set -e ;\
	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml up -d --build;\
	test_status_code=0 ;\
	docker build -t tests --network host -f tests/Dockerfile . || test_status_code=$$? ;\
	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down --rmi local --volumes --remove-orphans; \
    rm -f ${CFG_CALENDAR_FILE} ${CFG_SCHEDULER_FILE} ${CFG_SENDER_FILE}; \
	exit $$test_status_code ;

integration-test-cleanup: cfg-clean
	env `cat ./deployments/env.${ENV_NAME}` docker-compose -f ./deployments/docker-compose.yaml down --rmi local --volumes --remove-orphans;

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN_CALENDAR) -config ./configs/calendar_config.json
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.json
	$(BIN_SENDER) -config ./configs/sender_config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CALENDAR) version

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