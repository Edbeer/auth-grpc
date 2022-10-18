package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
	"github.com/Edbeer/microservices/internal/core"
	mockredis "github.com/Edbeer/microservices/internal/storage/redis/mock"
)

func TestService_CreateSession(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionStorage(ctrl)
	sessionService := newSessionService(mockSessionRedis)

	ctx := context.Background()
	session := &core.Session{}
	rT := "refresh token"

	mockSessionRedis.EXPECT().CreateSession(gomock.Any(), gomock.Eq(session), 10).Return(rT, nil)

	createdSession, err := sessionService.CreateSession(ctx, session, 10)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotEqual(t, createdSession, "")
}

func TestService_GetUserID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionStorage(ctrl)
	sessionService := newSessionService(mockSessionRedis)

	ctx := context.Background()
	session := &core.Session{}
	rT := "refresh token"

	mockSessionRedis.EXPECT().GetSessionByToken(gomock.Any(), gomock.Eq(rT)).Return(session, nil)

	session, err := sessionService.GetSessionByToken(ctx, rT)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, session)
}

func TestService_DeleteSessionByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRedis := mockredis.NewMockSessionStorage(ctrl)
	sessionService := newSessionService(mockSessionRedis)

	ctx := context.Background()
	rT := "refresh token"

	mockSessionRedis.EXPECT().DeleteSession(gomock.Any(), gomock.Eq(rT)).Return(nil)

	err := sessionService.DeleteSession(ctx, rT)
	require.NoError(t, err)
	require.Nil(t, err)
}