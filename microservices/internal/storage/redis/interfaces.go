//go:generate mockgen -source interfaces.go -destination mock/redis_storage_mock.go -package mock
package redstorage

import (
	"context"

	"guthub.com/Edbeer/microservices/internal/core"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, session *core.Session, expire int) (string, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}

type AccountStorage interface {
	GetByIDCtx(ctx context.Context, key string) (*core.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *core.User) error
}