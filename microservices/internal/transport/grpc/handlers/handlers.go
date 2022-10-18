package handlers

import "github.com/Edbeer/microservices/internal/config"

type Deps struct {
	AccountService AccountService
	SessionService SessionService
	Config         *config.Config
}

type Handlers struct {
	Account *accountHandler
	Example *Example
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		Account: newAccountHandler(deps.AccountService, deps.SessionService, deps.Config),
		Example: NewExample(),
	}
}
