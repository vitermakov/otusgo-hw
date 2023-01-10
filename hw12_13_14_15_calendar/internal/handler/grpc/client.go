package grpc

import (
	"fmt"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewSupportClient(apiAddr string) (events.SupportClient, error) {
	conn, err := grpc.Dial(apiAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("can't dial GRPC server: %w", err)
	}
	return events.NewSupportClient(conn), nil
}
