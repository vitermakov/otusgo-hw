# Environment
FROM golang:1.19 as build

ENV BIN_FILE /app/scheduler
ENV CODE_DIR /app

RUN mkdir -p ${CODE_DIR}
WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}
# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -a -o ${BIN_FILE} cmd/scheduler/*

# Release
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="scheduler"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build ${BIN_FILE} /sbin/scheduler

ENV CONFIG_FILE /etc/scheduler/config.json
COPY ./deployments/configs/scheduler_config.json ${CONFIG_FILE}

ENTRYPOINT ["/sbin/scheduler"]
