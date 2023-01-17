FROM golang:1.17

ENV DATA_DIR /data

RUN mkdir -p ${DATA_DIR}
WORKDIR ${DATA_DIR}

COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD ["/go/bin/goose", "-dir=migrations", "up"]
