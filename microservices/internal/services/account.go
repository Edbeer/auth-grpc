package service

import (
	"context"

	"github.com/google/uuid"
	"guthub.com/Edbeer/microservices/internal/core"
	"guthub.com/Edbeer/microservices/pkg/jwt"
	"guthub.com/Edbeer/microservices/pkg/utils"
)

// Token Manager interface
type Manager interface {
	GenerateJWTToken(user *core.User) (string, error)
	Parse(accessToken string) (*jwt.Claims, error)
}

type accountStorage interface {
	Create(ctx context.Context, user *core.User) (*core.User, error)
	FindByEmail(ctx context.Context, user *core.User) (*core.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*core.User, error)
}

type accountService struct {
	psql         accountStorage
	tokenManager Manager
}

func newAccountService(psql accountStorage, tokenManager Manager) *accountService {
	return &accountService{
		psql:         psql,
		tokenManager: tokenManager,
	}
}

func (a *accountService) SignUp(ctx context.Context, req *core.User) (*core.User, error) {
	foundUser, err := a.psql.FindByEmail(ctx, req)
	if err == nil || foundUser != nil {
		return nil, err
	}

	if err := req.PrepareCreate(); err != nil {
		return nil, err
	}

	if err := utils.ValidateStruct(ctx, req); err != nil {
		return nil, err
	}

	user, err := a.psql.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *accountService) SignIn(ctx context.Context, user *core.User) (*core.Token, error) {
	foundUser, err := a.psql.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenManager.GenerateJWTToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &core.Token{
		AccessToken: accessToken,
	}, nil
}
