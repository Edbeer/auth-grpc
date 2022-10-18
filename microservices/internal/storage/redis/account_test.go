package redstorage

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/Edbeer/microservices/internal/core"
)

func SetupAuthRedis() *accountStorage {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	authRedisStorage := newAccountStorage(client)
	return authRedisStorage 
}

func TestRedis_SetUserCtx(t *testing.T) {
	t.Parallel()

	authRedisStorage := SetupAuthRedis()

	t.Run("SetUserCtx", func(t *testing.T) {
		key := uuid.New().String()
		userId := uuid.New()
		u := &core.User{
			Uuid: userId,
			Name: "Pavel",
			Email: "edbeermtn@gmail.com",
		}

		err := authRedisStorage.SetUserCtx(context.Background(), key, 10, u)
		require.NoError(t, err)
		require.Nil(t, err)	
	})
}

func TestRedis_GetByIDCtx(t *testing.T) {
	t.Parallel()

	authRedisStorage := SetupAuthRedis()

	t.Run("GetByIDCtx", func(t *testing.T) {
		key := uuid.New().String()
		userId := uuid.New()
		u := &core.User{
			Uuid: userId,
			Name: "Pavel",
			Email: "edbeermtn@gmail.com",
		}

		err := authRedisStorage.SetUserCtx(context.Background(), key, 10, u)
		require.NoError(t, err)
		require.Nil(t, err)

		user, err := authRedisStorage.GetByIDCtx(context.Background(), key)
		require.NoError(t, err)
		require.NotNil(t, user)
	})
}