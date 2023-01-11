package queue

import (
	"errors"
	"fmt"

	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue/rabbit"
)

var ErrUnknownMPQType = errors.New("unsupported ampq type")

func NewRMPQConn(config common.Conn, logger logger.Logger, queueName string) *rabbit.MQConnection {
	MPQCfg := rabbit.Config{
		User:         config.User,
		Password:     config.Password,
		Host:         config.Host,
		Port:         config.Port,
		ExchangeName: fmt.Sprintf("%s_ex", queueName),
		ExchangeType: "direct",
		BindingKey:   fmt.Sprintf("%s_key", queueName),
	}
	return rabbit.New(MPQCfg, logger)
}

func NewProducer(config common.Queue, logger logger.Logger, queueName string) (queue.Producer, error) {
	if config.Type == "rabbitMq" {
		conn := NewRMPQConn(config.RabbitMQ, logger, queueName)
		return rabbit.NewProducer(conn)
	}
	return nil, ErrUnknownMPQType
}

func NewConsumer(config common.Queue, logger logger.Logger, queueName string) (queue.Consumer, error) {
	if config.Type == "rabbitMq" {
		conn := NewRMPQConn(config.RabbitMQ, logger, queueName)
		return rabbit.NewConsumer(conn), nil
	}
	return nil, ErrUnknownMPQType
}
