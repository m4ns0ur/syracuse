package postgres

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/go-toschool/syracuse"
)

// CitizensService ...
type CitizensService struct {
	Store *sqlx.DB
}

// Get ...
func (cs *CitizensService) Get(ID string) (*syracuse.Citizen, error) {
	query := squirrel.Select("*").From("users").Where("id = ?", ID).Where("deleted_at is null")

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := cs.Store.QueryRowx(sql, args...)

	c := &syracuse.Citizen{}
	if err := row.StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Select ...
func (cs *CitizensService) Select() ([]*syracuse.Citizen, error) {
	return nil, nil
}

// Create ...
func (cs *CitizensService) Create(c *syracuse.Citizen) error {
	sql, args, err := squirrel.
		Insert("users").
		Columns("email", "fullname").
		Values(c.Email, c.Fullname).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := cs.Store.QueryRowx(sql, args...)
	if err := row.StructScan(c); err != nil {
		return err
	}

	return nil
}

// Update ...
func (cs *CitizensService) Update(*syracuse.Citizen) error {
	return nil
}

// Delete ...
func (cs *CitizensService) Delete(*syracuse.Citizen) error {
	return nil
}
