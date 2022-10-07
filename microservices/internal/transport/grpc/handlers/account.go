package handlers

import (
	"context"

	accountpb "github.com/Edbeer/proto/api/account/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"guthub.com/Edbeer/microservices/internal/core"
)

type AccountService interface {
	SignUp(ctx context.Context, user *core.User) (*core.User, error)
	SignIn(ctx context.Context, user *core.User) (*core.Token, error)
}

type accountHandler struct {
	accountpb.UnimplementedAccountServiceServer
	service AccountService
}

func newAccountHandler(service AccountService) *accountHandler {
	return &accountHandler{
		service: service,
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

	token, err := a.service.SignIn(ctx, user)
	if err != nil {
		return nil, err
	}

	return &accountpb.SignInResponse{
		AccessToken: token.AccessToken,
	}, nil
}

func (a *accountHandler) RefreshTokens(ctx context.Context, req *accountpb.RefreshTokensRequest) (*accountpb.RefreshTokensResponse, error) {
	return nil, nil
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

