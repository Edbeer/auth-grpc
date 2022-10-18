package service

import (
	"context"

	"github.com/Edbeer/microservices/internal/core"
)

type sessionStorage interface {
	CreateSession(ctx context.Context, session *core.Session, expire int) (string, error)
	GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}

type sessionService struct {
	storage sessionStorage
}

func newSessionService(storage sessionStorage) *sessionService {
	return &sessionService{storage: storage}
}

func (s *sessionService) CreateSession(ctx context.Context, session *core.Session, expire int) (string, error) {
	return s.storage.CreateSession(ctx, session, expire)
}

func (s *sessionService) GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error) {
	return s.storage.GetSessionByToken(ctx, refreshToken)
}

func (s *sessionService) DeleteSession(ctx context.Context, refreshToken string) error {
	return s.storage.DeleteSession(ctx, refreshToken)
}

