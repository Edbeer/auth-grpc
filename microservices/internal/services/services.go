package service

import (
	"github.com/Edbeer/microservices/internal/storage/psql"
	"github.com/Edbeer/microservices/internal/storage/redis"
)

type Deps struct {
	Psql    *postgres.Storage
	Redis   *redstorage.Storage
	Manager Manager
}

type Service struct {
	Account *accountService
	Session *sessionService
}

func NewService(deps Deps) *Service {
	return &Service{
		Account: newAccountService(deps.Psql.Account, deps.Redis.Account, deps.Manager),
		Session: newSessionService(deps.Redis.Session),
	}
}
