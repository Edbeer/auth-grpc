package handlers

type Deps struct {
	AccountService AccountService
	SessionService SessionService
}

type Handlers struct {
	Account *accountHandler
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		Account: newAccountHandler(deps.AccountService, deps.SessionService),
	}
}