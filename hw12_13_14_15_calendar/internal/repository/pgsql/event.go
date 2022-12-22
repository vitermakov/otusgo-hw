package pgsql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
	_, err := stmt.ExecAndClose(ctx, er.pool)
	return err
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
	stmt := sqlf.From("events").Select("*")
	er.applySearch(stmt, search)
	stmt.Select("(select row_to_json(u) from users as u where events.owner_id=u.uuid)")

	events := make([]model.Event, 0)
	rows, err := er.pool.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		event, err := er.prepareModel(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (er EventRepo) prepareModel(row *sql.Rows) (model.Event, error) {
	var (
		id          sql.NullString
		ownerID     sql.NullString
		duration    sql.NullInt64
		description sql.NullString
		notifyTerm  sql.NullInt64
		userJSON    sql.NullString
		event       model.Event
	)
	if err := row.Scan(
		&id, &event.Title, &event.Date, &duration, &ownerID, &description,
		&notifyTerm, &event.CreatedAt, &event.UpdatedAt, &userJSON); err != nil {
		if err != nil {
			return event, err
		}
	}
	if id.Valid {
		guid, err := uuid.Parse(id.String)
		if err != nil {
			return event, err
		}
		event.ID = guid
	}
	if userJSON.Valid {
		var dtoUser struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		err := json.Unmarshal([]byte(userJSON.String), &dtoUser)
		if err != nil {
			return event, fmt.Errorf("error reading event owner: %w", err)
		}
		event.Owner = &model.User{
			Name:  dtoUser.Name,
			Email: dtoUser.Email,
		}
		guid, err := uuid.Parse(dtoUser.ID)
		if err != nil {
			return event, fmt.Errorf("error reading event owner id: %w", err)
		}
		event.Owner.ID = guid
	}
	if duration.Valid {
		event.NotifyTerm = time.Duration(duration.Int64)
	}
	if description.Valid {
		event.Description = description.String
	}
	if notifyTerm.Valid {
		event.NotifyTerm = time.Duration(notifyTerm.Int64)
	}
	return event, nil
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
