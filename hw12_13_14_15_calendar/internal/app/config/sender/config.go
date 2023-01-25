package sender

import (
	"errors"
	"fmt"
	"log"

	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
)

var ErrEmptyQueueListen = errors.New("empty queue name for sender listening")

type Config struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Logger      Logger `json:"logger"`
	APILogin    string `json:"apiLogin"`
	API         struct {
		Calendar common.API `json:"calendar"`
	} `json:"api"`
	Mailer common.Mailer `json:"mailer"`
	AMQP   common.Queue  `json:"amqp"`
	Notify Notify        `json:"notify"`
}

type Logger struct {
	FileName string `json:"fileName"`
	Level    string `json:"level"`
}

type Notify struct {
	QueueListen string `json:"queueListen"`
}

func New(fileName string) (Config, error) {
	var cfg Config
	if err := common.New(fileName, &cfg); err != nil {
		return cfg, fmt.Errorf("error reading configuaration from '%s': %w", fileName, err)
	}
	if len(cfg.Notify.QueueListen) == 0 {
		err := ErrEmptyQueueListen
		log.Println(err.Error())
		return Config{}, err
	}
	return cfg, nil
}
