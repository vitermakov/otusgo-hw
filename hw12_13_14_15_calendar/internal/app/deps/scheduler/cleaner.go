package scheduler

import (
	"context"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/grpc/pb/events"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"google.golang.org/protobuf/types/known/durationpb"
)

const defStoreTime = time.Hour * 24 * 365

type Cleaner struct {
	supportAPI events.SupportClient
	logger     logger.Logger

	storeTime time.Duration
}

func NewCleaner(
	api events.SupportClient, logger logger.Logger, storeTime string,
) *Cleaner {
	st, err := time.ParseDuration(storeTime)
	if err != nil {
		st = defStoreTime
		logger.Warn("wrong storeTime config value '%s', set default '%s'", storeTime, st.String())
	}
	return &Cleaner{supportAPI: api, logger: logger, storeTime: st}
}

func (cs Cleaner) DoAction(ctx context.Context) {
	_, err := cs.supportAPI.CleanupOldEvents(ctx, &events.CleanupReq{
		StoreTime: durationpb.New(cs.storeTime),
	})
	if err != nil {
		cs.logger.Error("cleaner action: error %s", err.Error())
	} else {
		cs.logger.Info("cleaner action: OK")
	}
}
