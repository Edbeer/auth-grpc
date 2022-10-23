package redstorage

import (
	"context"
	"encoding/json"

	"time"

	"github.com/Edbeer/microservices/internal/core"
	"github.com/go-redis/redis/v9"
	"github.com/opentracing/opentracing-go"
)

// Auth Storage
type accountStorage struct {
	redis *redis.Client
}

// Redis Account Storage constructor
func newAccountStorage(redis *redis.Client) *accountStorage {
	return &accountStorage{redis: redis}
}

// Get user by id
func (s *accountStorage) GetByIDCtx(ctx context.Context, key string) (*core.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRedis.GetByIDCtx")
	defer span.Finish()
	
	userBytes, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	user := &core.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Cache user with duration in seconds
func (s *accountStorage) SetUserCtx(ctx context.Context, key string, seconds int, user *core.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "accountRedis.SetUserCtx")
	defer span.Finish()
	
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := s.redis.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		return err
	}

	return nil
}
