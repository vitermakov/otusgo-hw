package pgsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/leporo/sqlf"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
)

type EventRepo struct {
	pool *sql.DB
}

func NewEventRepo(pool *sql.DB) repository.Event {
	return &EventRepo{pool: pool}
}

func (er EventRepo) Add(ctx context.Context, input model.EventCreate) (*model.Event, error) {
	guid := uuid.New()
	stmt := sqlf.InsertInto("events").
		Set("id", guid.String()).
		Set("title", input.Title).
		Set("date", input.Date).
		Set("duration", time.Duration(input.Duration)*time.Minute).
		Set("description", input.Description).
		Set("notify_term", input.NotifyTerm)
	// Returning("uuid").To(&guid)
	if input.OwnerID.ID() > 0 {
		stmt.Set("owner_id", input.OwnerID.String())
	}
	if input.Description != nil {
		stmt.Set("description", *input.Description)
	}
	if input.NotifyTerm != nil {
		stmt.Set("notify_term", time.Duration(*input.NotifyTerm)*time.Hour*24)
	}
	err := stmt.QueryRowAndClose(ctx, er.pool)
	if err != nil {
		return nil, err
	}
	events, err := er.GetList(ctx, model.EventSearch{ID: &guid})
	if err != nil {
		return nil, err
	}
	return &events[0], nil
}

func (er EventRepo) Update(ctx context.Context, input model.EventUpdate, search model.EventSearch) error {
	stmt := sqlf.Update("events").
		Set("updated_at", time.Now())
	er.applySearch(stmt, search)
	if input.Title != nil {
		stmt.Set("title", *input.Title)
	}
	if input.Date != nil {
		stmt.Set("date", *input.Date)
	}
	if input.Duration != nil {
		stmt.Set("duration", time.Duration(*input.Duration)*time.Minute)
	}
	if input.Description != nil {
		stmt.Set("description", *input.Description)
	}
	if input.NotifyTerm != nil {
		stmt.Set("notify_term", time.Duration(*input.NotifyTerm)*time.Hour*24)
	}
	if _, err := stmt.ExecAndClose(ctx, er.pool); err != nil {
		return err
	}
	return nil
}

func (er EventRepo) Delete(ctx context.Context, search model.EventSearch) error {
	stmt := sqlf.DeleteFrom("events")
	er.applySearch(stmt, search)
	if _, err := stmt.ExecAndClose(ctx, er.pool); err != nil {
		return err
	}
	return nil
}

// GetList не учитываем пагинацию, сортировку.
func (er EventRepo) GetList(ctx context.Context, search model.EventSearch) ([]model.Event, error) {
	var userJSON sql.NullString
	var dto struct {
		ID          string         `db:"id"`
		Title       string         `db:"title"`
		Date        time.Time      `db:"date"`
		Duration    int64          `db:"duration"`
		Description sql.NullString `db:"description"`
		NotifyTerm  sql.NullInt64  `db:"notify_term"`
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   time.Time      `db:"updated_at"`
	}
	stmt := sqlf.From("events").Bind(&dto)
	er.applySearch(stmt, search)
	stmt.Select("(select row_to_json(u) from users as u where events.owner_id=u.uuid)").To(&userJSON)

	events := make([]model.Event, 0)
	err := stmt.QueryAndClose(ctx, er.pool, func(row *sql.Rows) {
		event := model.Event{
			Title:     dto.Title,
			Date:      dto.Date,
			Duration:  time.Duration(dto.Duration * int64(time.Minute)),
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		}
		event.ID, _ = uuid.Parse(dto.ID)
		if userJSON.Valid {
			var dtoUser struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Email string `json:"email"`
			}
			_ = json.Unmarshal([]byte(userJSON.String), &dtoUser)
			event.Owner = &model.User{
				Name:  dtoUser.Name,
				Email: dtoUser.Email,
			}
			event.Owner.ID, _ = uuid.Parse(dtoUser.ID)
		}
		if dto.Description.Valid {
			event.Description = dto.Description.String
		}
		if dto.NotifyTerm.Valid {
			event.NotifyTerm = time.Duration(dto.NotifyTerm.Int64)
		}
		events = append(events, event)
	})
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (er EventRepo) applySearch(stmt *sqlf.Stmt, search model.EventSearch) {
	if search.ID != nil {
		stmt.Where("events.id = ?", search.ID.String())
	}
	if search.NotID != nil {
		stmt.Where("events.id != ?", search.ID.String())
	}
	if search.OwnerID != nil {
		stmt.Where("events.owner_id != ?", search.OwnerID.String())
	}
	if search.DateRange != nil {
		if search.TacDuration {
			stmt.Where("events.date + events.duration > ?", search.DateRange.GetFrom())
		} else {
			stmt.Where("events.date > ?", search.DateRange.GetFrom())
		}
		stmt.Where("events.date < ?", search.DateRange.GetTo())
	}
}
