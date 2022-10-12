package handlers

type Deps struct {
	AccountService AccountService
	SessionService SessionService
}

type Handlers struct {
	Account *accountHandler
	Example *Example
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		Account: newAccountHandler(deps.AccountService, deps.SessionService),
		Example: NewExample(),
	}
}