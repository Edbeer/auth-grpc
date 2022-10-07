package service

import "guthub.com/Edbeer/microservices/internal/storage/psql"

type Deps struct {
	Psql *postgres.Storage
	Manager Manager
}

type Service struct {
	Account *accountService
}

func NewService(deps Deps) *Service {
	return &Service{
		Account: newAccountService(deps.Psql.Account, deps.Manager),
	}
}