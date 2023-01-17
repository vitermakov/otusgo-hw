package rabbit

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

// Consumer чтение из очереди.
type Consumer struct {
	*MQConnection
}

func NewConsumer(conn *MQConnection) *Consumer {
	return &Consumer{MQConnection: conn}
}

type Worker func(context.Context, <-chan amqp.Delivery)

func (c *Consumer) Consume(ctx context.Context, queueName string) (<-chan queue.Message, error) {
	var err error
	if err = c.Connect(); err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	delivery, err := c.announceQueue(queueName)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	messages := make(chan queue.Message)
	go func() {
		var err error
		defer func() {
			close(messages)
			c.logger.Info("close messages channel")
		}()
		for {
			select {
			case del := <-delivery:
				if err = del.Ack(false); err != nil {
					c.logger.Error("error sending ack: %s", err.Error())
				}
				message := queue.Message(del.Body)
				select {
				case <-ctx.Done():
					return
				case messages <- message:
				}
			case <-c.done:
				delivery, err = c.reConnect(ctx, queueName)
				if err != nil {
					c.logger.Error("error reconnecting: %s", err.Error())
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return messages, nil
}
