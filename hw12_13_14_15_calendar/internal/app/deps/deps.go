package deps

import (
	"database/sql"
	"fmt"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/app/config"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository/memory"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository/pgsql"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers"
)

// Resources ресурсы внешних систем, таких как СУБД, очереди.
type Resources struct {
	DBPool *sql.DB
	// MQConn net.Conn
}

// Repos регистр репозиториев.
type Repos struct {
	Event repository.Event
	User  repository.User
}

func NewRepos(store config.Storage, res *Resources) (*Repos, error) {
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
			Event: pgsql.NewEventRepo(res.DBPool),
			User:  pgsql.NewUserRepo(res.DBPool),
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
}

// Services регистр сервисов.
type Services struct {
	Event  service.Event
	User   service.User
	Logger logger.Logger
	Auth   servers.AuthService
}

func NewServices(deps *Deps) *Services {
	repo := deps.Repos
	userServ := service.NewUserService(repo.User, deps.Logger)

	return &Services{
		Event:  service.NewEventService(repo.Event, deps.Logger, userServ),
		User:   userServ,
		Logger: deps.Logger,
		Auth:   service.NewAuthService(userServ),
	}
}
