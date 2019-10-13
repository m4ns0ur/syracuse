package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-toschool/syracuse"
	"github.com/jmoiron/sqlx"
)

// CitizensStore postgres implementation
type CitizensStore struct {
	Store *sqlx.DB
}

// Get gets a record from db
func (cs *CitizensStore) Get(q *syracuse.CitizensQuery) (*syracuse.Citizen, error) {
	query := squirrel.Select("*").From("users").Where("deleted_at is null")

	if q.ID == "" && q.Email == "" && q.FullName == "" {
		return nil, errors.New("must provide a query")
	}

	if q.ID != "" {
		query = query.Where("id = ?", q.ID)
	}

	if q.Email != "" {
		query = query.Where("email = ?", q.Email)
	}

	if q.FullName != "" {
		query = query.Where("full_name = ?", q.FullName)
	}
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

// Select returns a collection of users from db.
func (cs *CitizensStore) Select() ([]*syracuse.Citizen, error) {
	query := squirrel.Select("*").From("users").Where("deleted_at is null")

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cs.Store.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}

	cc := make([]*syracuse.Citizen, 0)

	for rows.Next() {
		c := &syracuse.Citizen{}
		if err := rows.StructScan(c); err != nil {
			return nil, err
		}
		cc = append(cc, c)
	}

	return cc, nil
}

// Create creates a new user.
func (cs *CitizensStore) Create(c *syracuse.Citizen) error {
	sql, args, err := squirrel.
		Insert("users").
		Columns("email", "full_name").
		Values(c.Email, c.FullName).
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

// Update updates the given user.
func (cs *CitizensStore) Update(c *syracuse.Citizen) error {
	sql, args, err := squirrel.Update("users").
		Set("email", c.Email).
		Set("full_name", c.FullName).
		Where("id = ?", c.ID).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	row := cs.Store.QueryRowx(sql, args...)
	return row.StructScan(c)
}

// Delete logical delete.
func (cs *CitizensStore) Delete(c *syracuse.Citizen) error {
	row := cs.Store.QueryRowx(
		"update users set deleted_at = $1 where id = $2 returning *",
		time.Now(), c.ID,
	)

	if err := row.StructScan(c); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return nil
}
