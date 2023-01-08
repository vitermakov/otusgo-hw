package calendar

import (
	"database/sql"
	"fmt"
	"github.com/benbjohnson/clock"
	common "github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository/pgsql"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
)

// Repos регистр репозиториев.
type Repos struct {
	Event repository.Event
	User  repository.User
}

func NewRepos(store common.Storage, dbPool *sql.DB) (*Repos, error) {
	var (
		repos *Repos
		err   error
	)
	switch store.Type {
	case "memory":
		repos = &Repos{
			Event: memory.NewEventRepo(),
			User:  memory.NewUserRepo(),
		}
	case "pgsql":
		repos = &Repos{
			Event: pgsql.NewEventRepo(dbPool),
			User:  pgsql.NewUserRepo(dbPool),
		}
	default:
		err = fmt.Errorf("unknown storage type '%s", store.Type)
	}
	return repos, err
}

// Deps зависимости.
type Deps struct {
	Repos  *Repos
	Logger logger.Logger
	Clock  clock.Clock
}

// Services регистр сервисов.
type Services struct {
	EventCRUD   service.EventCRUD
	EventNotify service.EventNotify
	EventClean  service.EventClean
	User        service.User
	Logger      logger.Logger
	Auth        servers.AuthService
}

func NewServices(deps *Deps) *Services {
	repo := deps.Repos
	userServ := service.NewUserService(repo.User, deps.Logger)

	return &Services{
		EventCRUD:   service.NewEventCRUDService(repo.Event, deps.Logger, userServ),
		EventNotify: service.NewEventNotifyService(repo.Event, deps.Logger, deps.Clock),
		EventClean:  service.NewEventCleanService(repo.Event, deps.Logger, deps.Clock),
		User:        userServ,
		Logger:      deps.Logger,
		Auth:        service.NewAuthService(userServ),
	}
}
