package queue

import (
	"context"
	"errors"
	"fmt"

	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/closer"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue/rabbit"
)

var ErrUnknownMPQType = errors.New("unsupported amqp type")

func NewAMQPConn(config common.Conn, logger logger.Logger, queueName string) *rabbit.MQConnection {
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

func NewProducer(
	config common.Queue, logger logger.Logger, queueName string,
) (queue.Producer, closer.CloseFunc, error) {
	if config.Type == "rabbitMq" {
		conn := NewAMQPConn(config.RabbitMQ, logger, queueName)
		producer, err := rabbit.NewProducer(conn)
		if err != nil {
			return nil, nil, err
		}
		return producer, connRabbitCloser(conn), nil
	}
	return nil, nil, ErrUnknownMPQType
}

func NewConsumer(
	config common.Queue, logger logger.Logger, queueName string,
) (queue.Consumer, closer.CloseFunc, error) {
	if config.Type == "rabbitMq" {
		conn := NewAMQPConn(config.RabbitMQ, logger, queueName)
		return rabbit.NewConsumer(conn), connRabbitCloser(conn), nil
	}
	return nil, nil, ErrUnknownMPQType
}

func connRabbitCloser(conn *rabbit.MQConnection) closer.CloseFunc {
	return func(ctx context.Context) error {
		return conn.Disconnect()
	}
}
