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
		Set("duration", fmt.Sprintf("%d seconds", int64(input.Duration.Seconds()))).
		Set("notify_status", model.NotifyStatusNone.String())
	if input.OwnerID.ID() > 0 {
		stmt.Set("owner_id", input.OwnerID.String())
	}
	if input.Description != nil {
		stmt.Set("description", *input.Description)
	}
	if input.NotifyTerm != nil {
		stmt.Set("notify_term", fmt.Sprintf("%d seconds", int64(input.NotifyTerm.Seconds())))
	}
	_, err := stmt.ExecAndClose(ctx, er.pool)
	if err != nil {
		return nil, err
	}
	events, err := er.GetList(ctx, model.EventSearch{ID: &guid})
	if err != nil {
		return nil, err
	}
	return &events[0], nil
}

func (er EventRepo) Update(ctx context.Context, input model.EventUpdate, search model.EventSearch) (int64, error) {
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
		stmt.Set("duration", fmt.Sprintf("%d seconds", int(input.Duration.Seconds())))
	}
	if input.Description != nil {
		stmt.Set("description", *input.Description)
	}
	if input.NotifyTerm != nil {
		stmt.Set("notify_term", fmt.Sprintf("%d seconds", int(input.NotifyTerm.Seconds())))
	}
	if input.NotifyStatus != nil {
		stmt.Set("notify_status", input.NotifyStatus.String())
	}
	res, err := stmt.ExecAndClose(ctx, er.pool)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (er EventRepo) Delete(ctx context.Context, search model.EventSearch) (int64, error) {
	stmt := sqlf.DeleteFrom("events")
	er.applySearch(stmt, search)
	res, err := stmt.ExecAndClose(ctx, er.pool)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetList не учитываем пагинацию, сортировку.
func (er EventRepo) GetList(ctx context.Context, search model.EventSearch) ([]model.Event, error) {
	stmt := sqlf.From("events").
		Select(`id, title, date, 
			EXTRACT(EPOCH FROM duration)::int, 
			description, 
			EXTRACT(EPOCH FROM notify_term)::int, 
			notify_status, created_at, updated_at`,
		)
	er.applySearch(stmt, search)
	stmt.Select("(select row_to_json(users) from users where events.owner_id=users.id) as owner")
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
		id, description, userJSON, notifyStatus sql.NullString
		duration, notifyTerm                    sql.NullInt64
		event                                   model.Event
	)
	if err := row.Scan(
		&id, &event.Title, &event.Date, &duration, &description,
		&notifyTerm, &notifyStatus, &event.CreatedAt, &event.UpdatedAt, &userJSON); err != nil {
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
		event.Duration = time.Duration(duration.Int64) * time.Second
	}
	if description.Valid {
		event.Description = description.String
	}
	if notifyTerm.Valid {
		event.NotifyTerm = time.Duration(notifyTerm.Int64) * time.Second
	}
	if notifyStatus.Valid {
		nf, err := model.ParseNotifyStatus(notifyStatus.String)
		if err != nil {
			return event, fmt.Errorf("error reading event owner: %w", err)
		}
		event.NotifyStatus = nf
	}
	return event, nil
}

func (er EventRepo) applySearch(stmt *sqlf.Stmt, search model.EventSearch) {
	if search.ID != nil {
		stmt.Where("events.id = ?", search.ID.String())
	}
	if search.NotID != nil {
		stmt.Where("events.id != ?", search.NotID.String())
	}
	if search.OwnerID != nil {
		stmt.Where("events.owner_id = ?", search.OwnerID.String())
	}
	if search.DateRange != nil {
		if search.TacDuration {
			stmt.Where("events.date + events.duration > ?", search.DateRange.GetFrom())
		} else {
			stmt.Where("events.date > ?", search.DateRange.GetFrom())
		}
		stmt.Where("events.date < ?", search.DateRange.GetTo())
	}
	if search.DateLess != nil {
		stmt.Where("events.date < ?", *search.DateLess)
	}
	if search.NeedNotifyTerm != nil {
		stmt.Where("events.notify_term IS NOT NULL")
		stmt.Where("events.notify_status = ?", model.NotifyStatusNone.String())
		stmt.Where("events.date - events.notify_term < ?", *search.NeedNotifyTerm)
		stmt.Where("events.date > ?", *search.NeedNotifyTerm)
	}
}

func (er EventRepo) BlockEvents4Notify(ctx context.Context, now time.Time) ([]model.Event, error) {
	stmt := sqlf.From("events").
		Select(`id, title, date, 
			EXTRACT(EPOCH FROM duration)::int, 
			description, 
			EXTRACT(EPOCH FROM notify_term)::int, 
			notify_status, created_at, updated_at`,
		)
	er.applySearch(stmt, model.EventSearch{
		NeedNotifyTerm: &now,
	})
	stmt.Select("(select row_to_json(users) from users where events.owner_id=users.id) as owner")

	events := make([]model.Event, 0)

	tx, err := er.pool.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	query := stmt.String() + " FOR UPDATE"
	rows, err := tx.QueryContext(ctx, query, stmt.Args()...)
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
		query = "UPDATE events SET notify_status=$1 WHERE id=$2"
		_, err = tx.ExecContext(ctx, query, model.NotifyStatusBlocked.String(), event.ID.String())
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return events, nil
}
