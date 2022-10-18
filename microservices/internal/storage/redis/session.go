package redstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/Edbeer/microservices/internal/core"
)

type sessionStorage struct {
	redis *redis.Client
}

func newSessionStorage(redis *redis.Client) *sessionStorage {
	return &sessionStorage{redis: redis}
}

func (s *sessionStorage) CreateSession(ctx context.Context, session *core.Session, expire int) (string, error) {
	session.RefreshToken = newRefreshToken()

	sessionBytes, err := json.Marshal(&session)
	if err != nil {
		return "", err
	}
	if err := s.redis.Set(ctx, session.RefreshToken, sessionBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return "", err
	}

	return session.RefreshToken, nil
}

// Get user id from session
func (s *sessionStorage) GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error) {

	sessionBytes, err := s.redis.Get(ctx, refreshToken).Bytes()
	if err != nil {
		return nil , err
	}
	session := &core.Session{}
	if err = json.Unmarshal(sessionBytes, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Delete session cookie
func (s *sessionStorage) DeleteSession(ctx context.Context, refreshToken string) error {
	if err := s.redis.Del(ctx, refreshToken).Err(); err != nil {
		return err
	}
	return nil
}

func newRefreshToken() string {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", b)
}