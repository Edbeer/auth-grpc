package handlers

import (
	"context"

	accountpb "github.com/Edbeer/proto/api/account/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"guthub.com/Edbeer/microservices/internal/core"
	"github.com/google/uuid"
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
}

func newAccountHandler(service AccountService, session SessionService) *accountHandler {
	return &accountHandler{
		service: service,
		session: session,
	}
}

func (a *accountHandler) SignUp(ctx context.Context, req *accountpb.SignUpRequest) (*accountpb.SignUpResponse, error) {
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
	user := &core.User{
		Email: req.Email,
		Pass:  req.Password,
	}

	userWithToken, err := a.service.SignIn(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.session.CreateSession(ctx, &core.Session{
		Uuid: userWithToken.User.Uuid,
	}, 3600)

	return &accountpb.SignInResponse{
		AccessToken:  userWithToken.AccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *accountHandler) RefreshTokens(ctx context.Context, req *accountpb.RefreshTokensRequest) (*accountpb.RefreshTokensResponse, error) {

	session, err := a.session.GetSessionByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := a.service.GetUserByID(ctx, session.Uuid)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.session.CreateSession(ctx, &core.Session{
		Uuid: user.User.Uuid,
	}, 3600)
	if err != nil {
		return nil, err
	}

	return &accountpb.RefreshTokensResponse{
		AccessToken:  user.AccessToken,
		RefreshToken: refreshToken,
	}, nil
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

func (a *accountHandler) getRefreshTokenFromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata.FromIncomingContext")
	}

	refreshToken := md.Get("refresh_token")
	if refreshToken[0] == "" {
		return "", status.Errorf(codes.PermissionDenied, "md.Get refresh_token")
	}

	return refreshToken[0], nil
}