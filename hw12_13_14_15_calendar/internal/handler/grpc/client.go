package grpc

import (
	"context"
	"fmt"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthFn func(ctx context.Context) context.Context

func NewSupportClient(apiAddr, apiLogin string) (events.SupportClient, AuthFn, error) {
	conn, err := grpc.Dial(apiAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("can't dial GRPC server: %w", err)
	}
	return events.NewSupportClient(conn), func(ctx context.Context) context.Context {
		meta := metadata.New(nil)
		meta.Append("authorization", apiLogin)
		return metadata.NewOutgoingContext(ctx, meta)
	}, nil
}
