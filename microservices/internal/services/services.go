package services

import "guthub.com/Edbeer/microservices/internal/storage/psql"

type Deps struct {
	Psql *psql.Storage
}

type Service struct {
	Account *accountService
}

func NewService(deps Deps) *Service {
	return &Service{
		Account: newAccountService(deps.Psql.Account),
	}
}