# Environment
FROM golang:1.19 as build

ENV BIN_FILE /app/sender
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
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -a -o ${BIN_FILE} cmd/sender/*

# Release
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="sender"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build ${BIN_FILE} /sbin/sender

ARG CONFIG_FILE
ENV CONFIG_FILE_LOCAL /etc/sender/config.json
COPY ${CONFIG_FILE} ${CONFIG_FILE_LOCAL}

ENTRYPOINT ["/sbin/sender"]
