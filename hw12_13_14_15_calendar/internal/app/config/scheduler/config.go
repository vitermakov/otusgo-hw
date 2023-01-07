package scheduler

import (
	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
)

type Config struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Logger      Logger `json:"logger"`
	APIs        struct {
		Calendar common.API `json:"calendar"`
	} `json:"APIs"`
	MPQ     common.Queue `json:"mpq"`
	Cleanup Cleanup      `json:"cleanup"`
	Notify  Notify       `json:"notify"`
}

type Logger struct {
	FileName string `json:"fileName"`
	Level    string `json:"level"`
}

type Cleanup struct {
	TimeLive string `json:"timeLive"` // с единицей измерения: 1y
}

type Notify struct {
	CheckingTime string `json:"checkingTime"` // с единицей измерения: 1m
	QueuePublish string `json:"queuePublish"`
}
