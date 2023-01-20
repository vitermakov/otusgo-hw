package calendar

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest"
)

type Client interface {
	Create(context.Context, dto.EventCreate) (*dto.Event, error)
	Update(context.Context, string, dto.EventUpdate) error
	GetByID(context.Context, string) (*dto.Event, error)
	GetListOnDate(context.Context, string, time.Time) ([]dto.Event, error)
	Delete(context.Context, string) error
}

type Auth struct {
	email string
}

func NewAuth(email string) rest.ClientAuth {
	return &Auth{email: email}
}

func (a Auth) Authorize(request *http.Request) error {
	if len(a.email) == 0 {
		return errors.New("empty user auth email")
	}
	request.Header.Set("Authorization", a.email)
	return nil
}

type ClientImpl struct {
	api *rest.Client
}

func NewClient(baseURL, authEmail string) Client {
	return &ClientImpl{
		api: rest.NewClient(baseURL, rest.WithAuth(NewAuth(authEmail))),
	}
}

func (c ClientImpl) Create(ctx context.Context, input dto.EventCreate) (*dto.Event, error) {
	var event = new(dto.Event)
	resp, err := c.api.Post(ctx, "/events", input)
	if err != nil {
		return nil, err
	}
	if err = rest.EncodeResponse(resp, event, true); err != nil {
		return nil, err
	}
	return event, nil
}

func (c ClientImpl) Update(ctx context.Context, id string, input dto.EventUpdate) error {
	resp, err := c.api.Put(ctx, fmt.Sprintf("/events/%s", id), input)
	if err != nil {
		return err
	}
	return rest.EncodeResponse(resp, nil, true)
}

func (c ClientImpl) GetByID(ctx context.Context, id string) (*dto.Event, error) {
	var event = new(dto.Event)
	resp, err := c.api.Get(ctx, fmt.Sprintf("/events/%s", id), nil)
	if err != nil {
		return nil, err
	}
	if err = rest.EncodeResponse(resp, event, false); err != nil {
		return nil, err
	}
	return event, nil
}

func (c ClientImpl) Delete(ctx context.Context, id string) error {
	resp, err := c.api.Delete(ctx, fmt.Sprintf("/events/%s", id), nil)
	if err != nil {
		return err
	}
	return rest.EncodeResponse(resp, nil, true)
}

func (c ClientImpl) GetListOnDate(ctx context.Context, rangeType string, t time.Time) ([]dto.Event, error) {
	var events []dto.Event
	resp, err := c.api.Get(
		ctx,
		fmt.Sprintf("/events/list/%s", rangeType),
		map[string]interface{}{"date": t.Format(time.RFC3339)},
	)
	if err != nil {
		return nil, err
	}
	if err = rest.EncodeResponse(resp, &events, false); err != nil {
		return nil, err
	}
	return events, nil
}
