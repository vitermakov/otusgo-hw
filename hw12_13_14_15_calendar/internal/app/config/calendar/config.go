package calendar

import (
	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
)

type Config struct {
	ServiceID   string        `json:"serviceId"`
	ServiceName string        `json:"serviceName"`
	Logger      common.Logger `json:"logger"`
	Servers     struct {
		HTTP common.Server `json:"http"`
		GRPC common.Server `json:"grpc"`
	} `json:"servers"`
	Storage common.Storage `json:"storage"`
}
