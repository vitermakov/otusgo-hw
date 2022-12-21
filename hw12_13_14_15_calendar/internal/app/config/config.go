package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Logger      Logger `json:"logger"`
	Servers     struct {
		HTTP Server `json:"http"`
		GRPC Server `json:"grpc"`
	} `json:"servers"`
	Storage       Storage       `json:"storage"`
	BgParams      BgParams      `json:"bgParams"`
	Notifications Notifications `json:"notifications"`
}

type Logger struct {
	FileName string `json:"fileName"`
	Level    string `json:"level"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type BgParams struct {
	TimeLive int `json:"timeLive"`
}

type Storage struct {
	Type   string  `json:"type"`
	PGConn SQLConn `json:"pgsql"`
}

type Queue struct {
	Type     string `json:"type"`
	RabbitMQ Conn   `json:"rabbitMq"`
}

type Notifications struct {
	DefaultTerm int    `json:"defaultTerm"`
	QueueName   string `json:"queueName"`
}

type Conn struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SQLConn struct {
	Conn
	DBName string `json:"dbName"`
}

// New пока используем обычный encode/json.
func New(fileName string) (Config, error) {
	var config Config
	bs, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	if err = json.Unmarshal(bs, &config); err != nil {
		return Config{}, err
	}
	return config, err
}
