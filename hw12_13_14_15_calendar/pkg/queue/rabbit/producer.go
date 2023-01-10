package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/queue"
)

// Producer отправка сообщений в очередь.
type Producer struct {
	*MQConnection
}

func NewProducer(conn *MQConnection) (*Producer, error) {
	err := conn.connect()
	if err != nil {
		return nil, err
	}
	return &Producer{MQConnection: conn}, nil
}

func (p *Producer) Produce(message queue.Message) error {
	ch, err := p.conn.Channel()
	defer func() {
		if err := ch.Close(); err != nil {
			p.logger.Error("can't close AMQP connection")
		}
	}()
	if err != nil {
		return fmt.Errorf("can't get channel from AMQP connection: %w", err)
	}
	err = ch.Publish(
		p.config.ExchangeName, // exchange
		p.config.BindingKey,   // routing key
		true,                  // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json; charset=utf-8",
			Body:         message,
		})
	if err != nil {
		return fmt.Errorf("can't publish message into RabbitMQ: %w", err)
	}
	return nil
}
