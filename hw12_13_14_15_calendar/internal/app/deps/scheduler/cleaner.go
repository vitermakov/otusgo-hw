package scheduler

import (
	"context"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Cleaner struct {
	supportAPI events.SupportClient
	authAPI    grpc.AuthFn
	logger     logger.Logger

	storeTime time.Duration
}

func NewCleaner(
	api events.SupportClient, authAPI grpc.AuthFn, logger logger.Logger, storeTime time.Duration,
) *Cleaner {
	return &Cleaner{supportAPI: api, authAPI: authAPI, logger: logger, storeTime: storeTime}
}

func (cs Cleaner) DoAction(ctx context.Context) {
	_, err := cs.supportAPI.CleanupOldEvents(
		cs.authAPI(ctx),
		&events.CleanupReq{
			StoreTime: durationpb.New(cs.storeTime),
		},
	)
	if err != nil {
		cs.logger.Error("cleaner action: error %s", err.Error())
	} else {
		cs.logger.Info("cleaner action: OK")
	}
}
