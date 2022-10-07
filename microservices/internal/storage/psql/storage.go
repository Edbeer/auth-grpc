package postgres

import "github.com/jmoiron/sqlx"

type Storage struct {
	Account *accountStorage
}

func NewStorage(psql *sqlx.DB) *Storage {
	return &Storage{
		Account: newAccountStorage(psql),
	}
}