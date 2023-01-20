package scheduler

import (
	"errors"
	"fmt"
	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/utils/jsonx"
	"log"
)

var ErrEmptyQueuePublish = errors.New("empty queue name for notifier publishing")

const (
	defCleanupCheckingTime = "1d"
	defCleanupStoreTime    = "1y"
	defNotifyCheckingTime  = "1m"
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
	CheckingTime jsonx.Duration `json:"checkingTime"` // с единицей измерения: 1d
	StoreTime    jsonx.Duration `json:"storeTime"`    // с единицей измерения: 1y
}

type Notify struct {
	CheckingTime jsonx.Duration `json:"checkingTime"` // с единицей измерения: 1m
	QueuePublish string         `json:"queuePublish"`
}

func New(fileName string) (Config, error) {
	var cfg Config
	if err := common.New(fileName, &cfg); err != nil {
		return cfg, fmt.Errorf("error reading configuaration from '%s': %w", fileName, err)
	}
	if !cfg.Cleanup.CheckingTime.Valid() {
		log.Printf(
			"wrong checkingTime cleaner config value, set default '%s'\n", defCleanupCheckingTime,
		)
		cfg.Cleanup.CheckingTime, _ = jsonx.ParseDuration(defCleanupCheckingTime)
	}

	if !cfg.Cleanup.StoreTime.Valid() {
		log.Printf(
			"wrong storeTime cleaner config value, set default '%s'\n", defCleanupStoreTime,
		)
		cfg.Cleanup.StoreTime, _ = jsonx.ParseDuration(defCleanupStoreTime)
	}

	if !cfg.Notify.CheckingTime.Valid() {
		log.Printf(
			"wrong checkingTime notifier config value, set default '%s'\n", defNotifyCheckingTime,
		)
		cfg.Notify.CheckingTime, _ = jsonx.ParseDuration(defNotifyCheckingTime)
	}

	if len(cfg.Notify.QueuePublish) == 0 {
		err := ErrEmptyQueuePublish
		log.Println(err.Error())
		return Config{}, err
	}
	return cfg, nil
}
