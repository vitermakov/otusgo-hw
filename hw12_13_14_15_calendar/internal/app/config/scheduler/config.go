package scheduler

import (
	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
)

type Config struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Logger      Logger `json:"logger"`
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
	CheckingTime string `json:"checkingTime"` // с единицей измерения: 1d
	StoreTime    string `json:"storeTime"`    // с единицей измерения: 1y
}

type Notify struct {
	CheckingTime string `json:"checkingTime"` // с единицей измерения: 1m
	QueuePublish string `json:"queuePublish"`
}
