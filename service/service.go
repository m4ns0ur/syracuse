package service

import (
	"github.com/go-toschool/syracuse"
	"github.com/go-toschool/syracuse/database"
)

// Citizens postgres implementation
type Citizens struct {
	Store database.Store
}

// GetByID gets a record from db by id.
func (cs *Citizens) GetByID(id string) (*syracuse.Citizen, error) {
	q := &syracuse.CitizensQuery{
		ID: id,
	}
	return cs.Store.Get(q)
}

// GetByEmail gets a record from db by email.
func (cs *Citizens) GetByEmail(email string) (*syracuse.Citizen, error) {
	q := &syracuse.CitizensQuery{
		Email: email,
	}
	return cs.Store.Get(q)
}

// GetByFullname gets a record from db by fullname.
func (cs *Citizens) GetByFullname(fullname string) (*syracuse.Citizen, error) {
	q := &syracuse.CitizensQuery{
		FullName: fullname,
	}
	return cs.Store.Get(q)
}

// Select returns a collectio of users from db.
func (cs *Citizens) Select() ([]*syracuse.Citizen, error) {
	return cs.Store.Select()
}

// Create creates a new user.
func (cs *Citizens) Create(c *syracuse.Citizen) error {
	return cs.Store.Create(c)
}

// Update updates the given user.
func (cs *Citizens) Update(c *syracuse.Citizen) error {
	return cs.Store.Update(c)
}

// Delete logical delete.
func (cs *Citizens) Delete(c *syracuse.Citizen) error {
	return cs.Store.Delete(c)
}
