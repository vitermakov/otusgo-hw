package calendar

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/client"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/handler/http/dto"
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

func NewAuth(email string) client.IAuth {
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
	api *client.API
}

func NewClient(baseURL, authEmail string) Client {
	return &ClientImpl{
		api: client.NewAPI(baseURL, client.WithAuth(NewAuth(authEmail))),
	}
}

func (c ClientImpl) Create(ctx context.Context, input dto.EventCreate) (*dto.Event, error) {
	var event *dto.Event
	resp, err := c.api.Post(ctx, "/events", input)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = client.EncodeResponse(resp, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c ClientImpl) Update(ctx context.Context, id string, input dto.EventUpdate) error {
	resp, err := c.api.Put(ctx, fmt.Sprintf("/events/%s", id), input)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return client.EncodeResponse(resp, nil)
}

func (c ClientImpl) GetByID(ctx context.Context, id string) (*dto.Event, error) {
	var event *dto.Event
	resp, err := c.api.Get(ctx, fmt.Sprintf("/events/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = client.EncodeResponse(resp, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (c ClientImpl) GetListOnDate(ctx context.Context, rangeType string, t time.Time) ([]dto.Event, error) {
	var events []dto.Event
	resp, err := c.api.Get(
		ctx,
		fmt.Sprintf("/events/list/%s", rangeType),
		map[string]interface{}{"date": t.String()},
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = client.EncodeResponse(resp, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (c ClientImpl) Delete(ctx context.Context, id string) error {
	resp, err := c.api.Delete(ctx, fmt.Sprintf("/events/%s", id), nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return client.EncodeResponse(resp, nil)
}
