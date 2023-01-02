package grpc

import (
	"context"
	"fmt"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

type LoggerInterceptor struct {
	logger logger.Logger
}

func NewLoggerInterceptor(logger logger.Logger) *LoggerInterceptor {
	return &LoggerInterceptor{logger: logger}
}

func (i *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		timeStart := time.Now()
		meta, ok := metadata.FromIncomingContext(ctx)
		resp, err := handler(ctx, req)
		ua := ""
		if ok {
			ua = strings.Join(meta.Get("user-agent"), " ")
		}
		i.logger.Info(
			fmt.Sprintf(
				"Method: %s\tDuration: %s\tError: %v\tUser-Agent: \"%s\"", info.FullMethod, time.Since(timeStart).String(), err, ua,
			),
		)
		return resp, err
	}
}
