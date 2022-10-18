package redstorage

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/Edbeer/microservices/internal/core"
)

func SetupSessionRedis() *sessionStorage {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	sessionRedisStorage := newSessionStorage(client)
	return sessionRedisStorage
}

func TestRedis_CreateSession(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("CreateSession", func(t *testing.T) {
		refreshToken := newRefreshToken()
		session := &core.Session{
			RefreshToken: refreshToken,
		}

		s, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, s, "")
	})
}

func TestRedis_GetSessionByToken(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("GetSessionByToken", func(t *testing.T) {
		userId := uuid.New()
		refreshToken := newRefreshToken()
		session := &core.Session{
			RefreshToken: refreshToken,
			Uuid: userId,
		}

		createdSession, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, createdSession, "")

		s, err := sessionRedisStorage.GetSessionByToken(context.Background(), createdSession)
		require.NoError(t, err)
		require.NotNil(t, s)
	})
}

func TestRedis_DeleteSession(t *testing.T) {
	t.Parallel()

	sessionRedisStorage := SetupSessionRedis()

	t.Run("DeleteSession", func(t *testing.T) {
		userId := uuid.New()
		refreshToken := newRefreshToken()
		session := &core.Session{
			RefreshToken: refreshToken,
			Uuid: userId,
		}

		createdSession, err := sessionRedisStorage.CreateSession(context.Background(), session, 10)
		require.NoError(t, err)
		require.NotEqual(t, createdSession, "")

		err = sessionRedisStorage.DeleteSession(context.Background(), createdSession)
		require.NoError(t, err)
		require.Nil(t, err)
	})
}