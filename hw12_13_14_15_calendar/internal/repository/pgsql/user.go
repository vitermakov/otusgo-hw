package pgsql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/leporo/sqlf"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/repository"
)

type UserRepo struct {
	pool *sql.DB
}

func NewUserRepo(pool *sql.DB) repository.User {
	return &UserRepo{pool: pool}
}

func (ur UserRepo) Add(ctx context.Context, input model.UserCreate) (*model.User, error) {
	guid := uuid.New()
	stmt := sqlf.InsertInto("users").
		Set("id", guid.String()).
		Set("name", input.Name).
		Set("email", input.Email)
	err := stmt.QueryRowAndClose(ctx, ur.pool)
	if err != nil {
		return nil, err
	}
	users, err := ur.GetList(ctx, model.UserSearch{ID: &guid})
	if err != nil {
		return nil, err
	}
	return &users[0], nil
}

func (ur UserRepo) Update(ctx context.Context, input model.UserUpdate, search model.UserSearch) (int64, error) {
	stmt := sqlf.Update("users")
	ur.applySearch(stmt, search)
	if input.Name != nil {
		stmt.Set("name", *input.Name)
	}
	if input.Email != nil {
		stmt.Set("email", *input.Email)
	}
	res, err := stmt.ExecAndClose(ctx, ur.pool)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (ur UserRepo) Delete(ctx context.Context, search model.UserSearch) (int64, error) {
	stmt := sqlf.DeleteFrom("users")
	ur.applySearch(stmt, search)
	res, err := stmt.ExecAndClose(ctx, ur.pool)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// GetList не учитываем пагинацию, сортировку.
func (ur UserRepo) GetList(ctx context.Context, search model.UserSearch) ([]model.User, error) {
	stmt := sqlf.From("users").Select("*")
	ur.applySearch(stmt, search)
	users := make([]model.User, 0)
	rows, err := ur.pool.QueryContext(ctx, stmt.String(), stmt.Args()...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		user, err := ur.prepareModel(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (ur UserRepo) prepareModel(row *sql.Rows) (model.User, error) {
	var (
		id    sql.NullString
		name  sql.NullString
		email sql.NullString
		user  model.User
	)
	if err := row.Scan(&id, &name, &email); err != nil {
		if err != nil {
			return user, err
		}
	}
	if id.Valid {
		guid, err := uuid.Parse(id.String)
		if err != nil {
			return user, err
		}
		user.ID = guid
	}
	if name.Valid {
		user.Name = name.String
	}
	if email.Valid {
		user.Email = email.String
	}
	return user, nil
}

func (ur UserRepo) applySearch(stmt *sqlf.Stmt, search model.UserSearch) {
	if search.ID != nil {
		stmt.Where("users.id = ?", search.ID.String())
	}
	if search.Email != nil {
		stmt.Where("users.email = ?", *search.Email)
	}
}
