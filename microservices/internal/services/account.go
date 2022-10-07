package services

type accountStorage interface {}

type accountService struct {
	psql accountStorage
}

func newAccountService(psql accountStorage) *accountService { 
	return &accountService{psql: psql}
}