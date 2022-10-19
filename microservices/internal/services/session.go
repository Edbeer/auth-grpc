package service

import (
	"context"

	"github.com/Edbeer/microservices/internal/core"
	"github.com/opentracing/opentracing-go"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionService.CreateSession")
	defer span.Finish()
	return s.storage.CreateSession(ctx, session, expire)
}

func (s *sessionService) GetSessionByToken(ctx context.Context, refreshToken string) (*core.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionService.GetSessionByToken")
	defer span.Finish()
	return s.storage.GetSessionByToken(ctx, refreshToken)
}

func (s *sessionService) DeleteSession(ctx context.Context, refreshToken string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionService.DeleteSession")
	defer span.Finish()
	return s.storage.DeleteSession(ctx, refreshToken)
}

