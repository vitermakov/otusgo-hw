package sender

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
	Mailer common.Mailer `json:"mailer"`
	MPQ    common.Queue  `json:"mpq"`
	Notify Notify        `json:"notify"`
}

type Logger struct {
	FileName string `json:"fileName"`
	Level    string `json:"level"`
}

type Notify struct {
	QueueListen string `json:"queueListen"`
}
