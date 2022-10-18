package handlers

import (
	"context"
	"testing"


	accountpb "github.com/Edbeer/proto/api/account/v1"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/Edbeer/microservices/internal/config"
	"github.com/Edbeer/microservices/internal/core"
	"github.com/Edbeer/microservices/internal/services/mock"
)

func TestUsersService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	account := mock.NewMockAccountService(ctrl)
	session := mock.NewMockSessionService(ctrl)
	authServerGRPC := newAccountHandler(account, session, nil)

	reqValue := &accountpb.SignUpRequest{
		Email:    "edbeer123mtn@gmail.com",
		Name:     "PavelV",
		Password: "12345678",
		Role:     "user",
	}

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		user := &core.User{
			Uuid:  userID,
			Email: reqValue.Email,
			Name:  reqValue.Name,
			Pass:  reqValue.Password,
			Role:  reqValue.Role,
		}

		account.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(user, nil)

		response, err := authServerGRPC.SignUp(context.Background(), reqValue)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, reqValue.Email, response.User.Email)
	})
}

func TestUsersService_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	account := mock.NewMockAccountService(ctrl)
	session := mock.NewMockSessionService(ctrl)
	cfg := &config.Config{
		Session: config.SessionConfig{
			ExpireAt: 10,
		},
	}
	authServerGRPC := newAccountHandler(account, session, cfg)

	reqValue := &accountpb.SignInRequest{
		Email:    "edbeer123mtn@gmail.com",
		Password: "12345678",
	}

	t.Run("Login", func(t *testing.T) {
		t.Parallel()
		token := "refresh token"
		user := &core.User{
			Email: "edbeer123mtn@gmail.com",
			Pass:  "12345678",
		}
		userID := uuid.New()
		userT := &core.UserWithToken{
			User: &core.User{
				Uuid: userID,
			},
		}
		sess := &core.Session{
			Uuid: userID,
		}
		account.EXPECT().SignIn(context.Background(), gomock.Eq(user)).Return(userT, nil)
		session.EXPECT().CreateSession(context.Background(), gomock.Eq(sess), cfg.Session.ExpireAt).Return(token, nil)

		response, err := authServerGRPC.SignIn(context.Background(), reqValue)
		require.NoError(t, err)
		require.NotNil(t, response)
	})
}