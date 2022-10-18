package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"guthub.com/Edbeer/microservices/internal/core"
	"github.com/google/uuid"
)

type accountStorage struct {
	psql *sqlx.DB
}

func newAccountStorage(psql *sqlx.DB) *accountStorage {
	return &accountStorage{psql: psql}
}

func (a *accountStorage) Create(ctx context.Context, user *core.User) (*core.User, error) {
	u := &core.User{}
	query := `INSERT INTO users (name, email, password, role, created_at)
			VALUES ($1, $2, $3, $4, now())
			RETURNING *`
	if err := a.psql.QueryRowxContext(ctx, query,
		&user.Name, &user.Email, &user.Pass, &user.Role,
	).StructScan(u); err != nil {
		return nil, fmt.Errorf("UserStoragePsql.Create.StructScan: %v", err)
	}
	return u, nil
}

func (a *accountStorage) FindByEmail(ctx context.Context, user *core.User) (*core.User, error) {
	u := &core.User{}
	query := `SELECT user_id, name, email, password, role, created_at
		FROM users
		WHERE email = $1`
	if err := a.psql.QueryRowxContext(ctx, query, &user.Email).StructScan(u); err != nil {
		return nil, fmt.Errorf("UserStoragePsql.Create.StructScan: %v", err)
	}
	return u, nil
}

// Get user by id
func (a *accountStorage) GetUserByID(ctx context.Context, userID uuid.UUID) (*core.User, error) {
	
	u := &core.User{}
	query := `SELECT user_id, name, email, password, role, created_at
		FROM users
		WHERE user_id = $1`
	if err := a.psql.QueryRowxContext(ctx, query, userID).StructScan(u); err != nil {
		return nil, fmt.Errorf("AuthStoragePsql.GetUserByID.StructScan: %w", err)
	}

	return u, nil
}