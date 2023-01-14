package rabbit

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	backoffv3 "github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
)

type Config struct {
	User         string
	Password     string
	Host         string
	Port         int
	ExchangeName string
	ExchangeType string
	BindingKey   string
}

type MQConnection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	done    chan error
	config  Config
	logger  logger.Logger
}

func New(cfg Config, logger logger.Logger) *MQConnection {
	return &MQConnection{config: cfg, logger: logger}
}

func (r *MQConnection) Connect() error {
	var err error
	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(r.config.User, r.config.Password),
		Host:   net.JoinHostPort(r.config.Host, strconv.Itoa(r.config.Port)),
		Path:   "/",
	}
	r.conn, err = amqp.Dial(uri.String())
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	go func() {
		ncErr := <-r.conn.NotifyClose(make(chan *amqp.Error))
		r.logger.Error("channel closed: %s", ncErr.Error())
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		r.done <- fmt.Errorf("channel Closed: %w", ncErr)
	}()

	if err = r.channel.ExchangeDeclare(
		r.config.ExchangeName,
		r.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	return nil
}

func (r *MQConnection) Disconnect() error {
	if !r.conn.IsClosed() {
		return r.conn.Close()
	}
	return nil
}

// Задекларировать очередь, которую будем слушать.
func (r *MQConnection) announceQueue(queueName string) (<-chan amqp.Delivery, error) {
	queue, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %w", err)
	}
	// Число сообщений, которые можно подтвердить за раз.
	err = r.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %w", err)
	}
	// Создаём биндинг (правило маршрутизации).
	if err = r.channel.QueueBind(
		queue.Name,
		r.config.BindingKey,
		r.config.ExchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %w", err)
	}
	msgs, err := r.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %w", err)
	}

	return msgs, nil
}

func (r *MQConnection) reConnect(ctx context.Context, queue string) (<-chan amqp.Delivery, error) {
	be := backoffv3.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoffv3.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoffv3.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}
		<-time.After(d)
		if err := r.Connect(); err != nil {
			log.Printf("could not connect in reconnect call: %+v", err)
			continue
		}
		msgs, err := r.announceQueue(queue)
		if err != nil {
			fmt.Printf("Couldn't connect: %+v", err)
			continue
		}

		return msgs, nil
	}
}
