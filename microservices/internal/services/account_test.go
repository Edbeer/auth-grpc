package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"guthub.com/Edbeer/microservices/internal/config"
	"guthub.com/Edbeer/microservices/internal/core"
	psql "guthub.com/Edbeer/microservices/internal/storage/psql/mock"
	redis "guthub.com/Edbeer/microservices/internal/storage/redis/mock"

	"guthub.com/Edbeer/microservices/pkg/jwt"
)

func TestService_SignUp(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		GrpsServer: config.GrpcServerConfig{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.GrpsServer.JwtSecretKey)
	mockUserStorage := psql.NewMockAccountStorage(ctrl)
	userService := newAccountService(mockUserStorage, nil, manager)

	user := &core.User{
		Name:  "PavelV",
		Pass:  "12345678",
		Email: "edbeer123mtn@gmail.com",
		Role:  "user",
	}

	ctx := context.Background()
	// span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "AccountSevice.SignUp")
	// defer span.Finish()

	mockUserStorage.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockUserStorage.EXPECT().Create(ctx, gomock.Eq(user)).Return(user, nil)

	createdUser, err := userService.SignUp(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Nil(t, err)
}

func TestService_SignIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		GrpsServer: config.GrpcServerConfig{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.GrpsServer.JwtSecretKey)
	mockUserStorage := psql.NewMockAccountStorage(ctrl)
	mockAuthRedis := redis.NewMockAccountStorage(ctrl)
	userService := newAccountService(mockUserStorage, mockAuthRedis, manager)

	user := &core.User{
		Pass:  "12345678",
		Email: "edbeer123mtn@gmail.com",
	}

	ctx := context.Background()
	// span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "AccountSevice.SignIn")
	// defer span.Finish()

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &core.User{
		Pass:  string(hashPassword),
		Email: "edbeermtn@gmail.com",
	}

	mockUserStorage.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := userService.SignIn(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}

func TestService_GetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := &config.Config{
		GrpsServer: config.GrpcServerConfig{
			JwtSecretKey: "secret",
		},
	}

	manager, _ := jwt.NewManager(config.GrpsServer.JwtSecretKey)
	mockUserStorage := psql.NewMockAccountStorage(ctrl)
	mockAuthRedis := redis.NewMockAccountStorage(ctrl)
	userService := newAccountService(mockUserStorage, mockAuthRedis, manager)
	user := &core.User{
		Pass:  "12345678",
		Email: "edbeer123mtn@gmail.com",
	}

	ctx := context.Background()
	// span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "AccountSevice.GetUserByID")
	// defer span.Finish()

	mockAuthRedis.EXPECT().GetByIDCtx(ctx, user.Uuid.String()).Return(nil, nil)
	mockUserStorage.EXPECT().GetUserByID(ctx, gomock.Eq(user.Uuid)).Return(user, nil)
	mockAuthRedis.EXPECT().SetUserCtx(ctx, user.Uuid.String(), 3600, user).Return(nil)

	u, err := userService.GetUserByID(ctx, user.Uuid)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}
