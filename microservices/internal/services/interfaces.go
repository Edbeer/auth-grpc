//go:generate mockgen -source interfaces.go -destination mock/service_mock.go -package mock
package service

import (
	"context"

	"github.com/google/uuid"
	"guthub.com/Edbeer/microservices/internal/core"
)

type AccountService interface {
	SignUp(ctx context.Context, user *core.User) (*core.User, error)
	SignIn(ctx context.Context, user *core.User) (*core.UserWithToken, error)
	GetUserByID(ctx context.Context, userId uuid.UUID) (*core.UserWithToken, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, session *core.Session, expire int) (string, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}