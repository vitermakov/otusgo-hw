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
	// Returning("uuid").To(&guid)
	err := stmt.QueryRowAndClose(ctx, ur.pool)
	if err != nil {
		return nil, err
	}
	events, err := ur.GetList(ctx, model.UserSearch{ID: &guid})
	if err != nil {
		return nil, err
	}
	return &events[0], nil
}

func (ur UserRepo) Update(ctx context.Context, input model.UserUpdate, search model.UserSearch) error {
	stmt := sqlf.Update("users")
	ur.applySearch(stmt, search)
	if input.Name != nil {
		stmt.Set("name", *input.Name)
	}
	if input.Email != nil {
		stmt.Set("email", *input.Email)
	}
	if _, err := stmt.ExecAndClose(ctx, ur.pool); err != nil {
		return err
	}
	return nil
}

func (ur UserRepo) Delete(ctx context.Context, search model.UserSearch) error {
	stmt := sqlf.DeleteFrom("users")
	ur.applySearch(stmt, search)
	if _, err := stmt.ExecAndClose(ctx, ur.pool); err != nil {
		return err
	}
	return nil
}

// GetList не учитываем пагинацию, сортировку.
func (ur UserRepo) GetList(ctx context.Context, search model.UserSearch) ([]model.User, error) {
	var dto struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		Email string `db:"email"`
	}
	stmt := sqlf.From("users").Bind(&dto)
	ur.applySearch(stmt, search)
	users := make([]model.User, 0)
	err := stmt.QueryAndClose(ctx, ur.pool, func(row *sql.Rows) {
		user := model.User{
			Name:  dto.Name,
			Email: dto.Email,
		}
		user.ID, _ = uuid.Parse(dto.ID)
		users = append(users, user)
	})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur UserRepo) applySearch(stmt *sqlf.Stmt, search model.UserSearch) {
	if search.ID != nil {
		stmt.Where("users.id = ?", search.ID.String())
	}
	if search.Email != nil {
		stmt.Where("users.email = ?", *search.Email)
	}
}
