package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServiceId   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Logger      Logger `json:"logger"`
	Servers     struct {
		Http Server `json:"http"`
		Grpc Server `json:"grpc"`
	} `json:"servers"`
	Storage       Storage       `json:"storage"`
	BgParams      BgParams      `json:"bg_params"`
	Notifications Notifications `json:"notifications"`
}

type Logger struct {
	FileName string `json:"file_name"`
	Level    string `json:"level"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type BgParams struct {
	TimeLive int `json:"time_live"`
}

type Storage struct {
	Type   string  `json:"type"`
	PgConn SqlConn `json:"pgsql"`
}

type Queue struct {
	Type     string `json:"type"`
	RabbitMQ Conn   `json:"rabbit_mq"`
}

type Notifications struct {
	DefaultTerm int    `json:"default_term"`
	QueueName   string `json:"queue_name"`
}

type Conn struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SqlConn struct {
	Conn
	DbName string `json:"dbname"`
}

// New пока исользуем обычный encode/json
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
