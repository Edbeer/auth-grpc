package handlers

import (
	"context"

	"github.com/Edbeer/microservices/internal/config"
	"github.com/Edbeer/microservices/internal/core"
	accountpb "github.com/Edbeer/microservices/proto/api/account/v1"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type accountHandler struct {
	accountpb.UnimplementedAccountServiceServer
	service AccountService
	session SessionService
	config  *config.Config
}

func newAccountHandler(
	service AccountService,
	session SessionService,
	config *config.Config,
) *accountHandler {
	return &accountHandler{
		service: service,
		session: session,
		config:  config,
	}
}

func (a *accountHandler) SignUp(ctx context.Context, req *accountpb.SignUpRequest) (*accountpb.SignUpResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "account.SignUp")
	defer span.Finish()
	
	u := &core.User{
		Name:  req.Name,
		Email: req.Email,
		Pass:  req.Password,
		Role:  req.Role,
	}

	user, err := a.service.SignUp(ctx, u)
	if err != nil {
		return nil, err
	}

	return &accountpb.SignUpResponse{
		User: a.userToProto(user),
	}, nil
}

func (a *accountHandler) SignIn(ctx context.Context, req *accountpb.SignInRequest) (*accountpb.SignInResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "account.SignIn")
	defer span.Finish()
	
	user := &core.User{
		Email: req.Email,
		Pass:  req.Password,
	}

	userWithToken, err := a.service.SignIn(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.session.CreateSession(ctx, &core.Session{
		Uuid:     userWithToken.User.Uuid,
	}, a.config.Session.ExpireAt)

	return &accountpb.SignInResponse{
		AccessToken:  userWithToken.AccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *accountHandler) RefreshTokens(ctx context.Context, req *accountpb.RefreshTokensRequest) (*accountpb.RefreshTokensResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "account.RefreshTokens")
	defer span.Finish()

	session, err := a.session.GetSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.service.GetUserByID(ctx, session.Uuid)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.session.CreateSession(ctx, &core.Session{
		Uuid:     user.User.Uuid,
	}, a.config.Session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &accountpb.RefreshTokensResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout user, delete current session
func (a *accountHandler) Logout(ctx context.Context, req *accountpb.SignOutRequest) (*accountpb.SignOutResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "account.Logout")
	defer span.Finish()

	session, err := a.session.GetSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if err := a.session.DeleteSession(ctx, session.RefreshToken); err != nil {
		return nil, err
	}

	return &accountpb.SignOutResponse{}, nil
}

func (a *accountHandler) userToProto(user *core.User) *accountpb.User {
	return &accountpb.User{
		Uuid:      user.Uuid.String(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Pass,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
