# Environment
FROM golang:1.19 as build

ENV BIN_FILE /app/calendar
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
RUN CGO_ENABLED=0 go build -ldflags "$LDFLAGS" -a -o ${BIN_FILE} cmd/calendar/*

# Release
FROM alpine:latest

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar"
LABEL MAINTAINERS="student@otus.ru"

COPY --from=build ${BIN_FILE} /sbin/calendar

ENV CONFIG_FILE /etc/calendar/config.json
COPY ./deployments/configs/calendar_config.json ${CONFIG_FILE}

EXPOSE 8080
EXPOSE 8088

ENTRYPOINT ["/sbin/calendar"]
