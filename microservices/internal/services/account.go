package service

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/Edbeer/microservices/internal/core"
	"github.com/Edbeer/microservices/pkg/jwt"
	"github.com/Edbeer/microservices/pkg/utils"
)

// Token Manager interface
type Manager interface {
	GenerateJWTToken(user *core.User) (string, error)
	Parse(accessToken string) (*jwt.Claims, error)
}

// postgres storage
type accountStorage interface {
	Create(ctx context.Context, user *core.User) (*core.User, error)
	FindByEmail(ctx context.Context, user *core.User) (*core.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*core.User, error)
}

// redis storage
type redisStorage interface {
	GetByIDCtx(ctx context.Context, key string) (*core.User, error)
	SetUserCtx(ctx context.Context, key string, seconds int, user *core.User) error
}

type accountService struct {
	psql         accountStorage
	redis        redisStorage
	tokenManager Manager
}

func newAccountService(psql accountStorage, redis redisStorage, tokenManager Manager) *accountService {
	return &accountService{
		psql:         psql,
		redis:        redis,
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

func (a *accountService) SignIn(ctx context.Context, user *core.User) (*core.UserWithToken, error) {
	foundUser, err := a.psql.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenManager.GenerateJWTToken(foundUser)
	if err != nil {
		return nil, err
	}

	return &core.UserWithToken{
		User:        foundUser,
		AccessToken: accessToken,
	}, nil
}

// Get user by id
func (a *accountService) GetUserByID(ctx context.Context, userId uuid.UUID) (*core.UserWithToken, error) {
	// Looking for a cached user
	cachedUser, err := a.redis.GetByIDCtx(ctx, userId.String())
	if err != nil {
		log.Println("%w", err)
	}
	// if cachedUser exist
	if cachedUser != nil {
		accessToken, err := a.tokenManager.GenerateJWTToken(cachedUser)
		if err != nil {
			return nil, err
		}
		return 	&core.UserWithToken{
			User:        cachedUser,
			AccessToken: accessToken,
		}, nil
	}

	foundUser, err := a.psql.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenManager.GenerateJWTToken(foundUser)
	if err != nil {
		return nil, err
	}

	if err := a.redis.SetUserCtx(ctx, userId.String(), 3600, foundUser); err != nil {
		fmt.Println("AccountService.GetByID.SetUserCtx: %w", err)
	}

	return &core.UserWithToken{
		User:        foundUser,
		AccessToken: accessToken,
	}, nil
}
