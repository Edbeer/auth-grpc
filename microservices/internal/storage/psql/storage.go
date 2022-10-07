package psql

type Storage struct {
	Account *accountStorage
}

func NewStorage() *Storage {
	return &Storage{
		Account: newAccountStorage(),
	}
}