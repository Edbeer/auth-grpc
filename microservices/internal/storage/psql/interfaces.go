//go:generate mockgen -source interfaces.go -destination mock/postgres_storage_mock.go -package mock
package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/Edbeer/microservices/internal/core"
)

type AccountStorage interface {
	Create(ctx context.Context, user *core.User) (*core.User, error)
	FindByEmail(ctx context.Context, user *core.User) (*core.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*core.User, error)
}