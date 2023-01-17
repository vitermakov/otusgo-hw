package scheduler

import (
	"errors"
	"fmt"
	"log"
	"time"

	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
)

var ErrEmptyQueuePublish = errors.New("empty queue name for notifier publishing")

const (
	defCleanupCheckingTime = time.Hour * 24
	defCleanupStoreTime    = time.Hour * 24 * 365
	defNotifyCheckingTime  = time.Minute
)

type Config struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Logger      Logger `json:"logger"`
	APILogin    string `json:"apiLogin"`
	API         struct {
		Calendar common.API `json:"calendar"`
	} `json:"api"`
	MPQ     common.Queue `json:"mpq"`
	Cleanup Cleanup      `json:"cleanup"`
	Notify  Notify       `json:"notify"`
}

type Logger struct {
	FileName string `json:"fileName"`
	Level    string `json:"level"`
}

type Cleanup struct {
	CheckingTime time.Duration `json:"checkingTime"` // с единицей измерения: 1d
	StoreTime    time.Duration `json:"storeTime"`    // с единицей измерения: 1y
}

type Notify struct {
	CheckingTime time.Duration `json:"checkingTime"` // с единицей измерения: 1m
	QueuePublish string        `json:"queuePublish"`
}

func New(fileName string) (Config, error) {
	var cfg Config
	if err := common.New(fileName, &cfg); err != nil {
		return cfg, fmt.Errorf("error reading configuaration from '%s': %w", fileName, err)
	}
	if cfg.Cleanup.CheckingTime <= 0 {
		log.Printf(
			"wrong checkingTime cleaner config value '%s', set default '%s'\n",
			cfg.Cleanup.CheckingTime,
			defCleanupCheckingTime,
		)
		cfg.Cleanup.StoreTime = defCleanupCheckingTime
	}
	if cfg.Cleanup.StoreTime <= 0 {
		log.Printf(
			"wrong storeTime cleaner config value '%s', set default '%s'\n",
			cfg.Cleanup.StoreTime,
			defCleanupStoreTime,
		)
		cfg.Cleanup.StoreTime = defCleanupStoreTime
	}
	if cfg.Notify.CheckingTime <= 0 {
		log.Printf(
			"wrong checkingTime notifier config value '%s', set default '%s'\n",
			cfg.Notify.CheckingTime,
			defNotifyCheckingTime,
		)
		cfg.Notify.CheckingTime = defNotifyCheckingTime
	}
	if len(cfg.Notify.QueuePublish) == 0 {
		err := ErrEmptyQueuePublish
		log.Println(err.Error())
		return Config{}, err
	}
	return cfg, nil
}
