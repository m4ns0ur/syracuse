package database

import (
	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/database/postgres"
	"github.com/jmoiron/sqlx"
)

// Store define the behavior for a Store.
type Store interface {
	Get(*syracuse.CitizensQuery) (*syracuse.Citizen, error)
	Select() ([]*syracuse.Citizen, error)
	Create(*syracuse.Citizen) error
	Update(*syracuse.Citizen) error
	Delete(*syracuse.Citizen) error
}

// NewPostgres return a new citizens service.
func NewPostgres(dsn string) (Store, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &postgres.CitizensStore{
		Store: db,
	}, nil
}
