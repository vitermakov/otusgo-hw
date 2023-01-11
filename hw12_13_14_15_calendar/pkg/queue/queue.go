package queue

import (
	"context"
	"encoding/json"
)

// Producer отправка сообщений в очередь.
type Producer interface {
	Produce(message Message) error
}

type Consumer interface {
	Consume(ctx context.Context, queue string) (<-chan Message, error)
}

type Message []byte

func EncMessage(object interface{}) (Message, error) {
	bs, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (m Message) Decode(object interface{}) error {
	return json.Unmarshal(m, object)
}
