package handlers

type Deps struct {
	AccountService AccountService
}

type Handlers struct {
	Account *accountHandler
}

func NewHandlers(deps Deps) *Handlers {
	return &Handlers{
		Account: newAccountHandler(deps.AccountService),
	}
}