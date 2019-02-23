package syracuse

import (
	"time"

	"github.com/go-toschool/syracuse/citizens"
)

// Citizen represents a user.
type Citizen struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	FullName string `json:"full_name" db:"full_name"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// Citizens basic method that need to be implemented in order to operate over users.
type Citizens interface {
	CitizensGetter

	Select() ([]*Citizen, error)
	Create(*Citizen) error
	Update(*Citizen) error
	Delete(*Citizen) error
}

// CitizensGetter getter define behavior to get records.
type CitizensGetter interface {
	GetByID(id string) (*Citizen, error)
	GetByEmail(email string) (*Citizen, error)
	GetByFullname(fullname string) (*Citizen, error)
}

// CitizensQuery represents queries that helps to query users.
type CitizensQuery struct {
	ID       string
	Email    string
	FullName string
}

// ToProto ...
func (c *Citizen) ToProto() *citizens.Citizen {
	return &citizens.Citizen{
		Id:        c.ID,
		FullName:  c.FullName,
		Email:     c.Email,
		CreatedAt: c.CreatedAt.UnixNano(),
		UpdatedAt: c.UpdatedAt.UnixNano(),
	}
}

// FromProto ...
func (c *Citizen) FromProto(cc *citizens.Citizen) *Citizen {
	c.ID = cc.Id
	c.FullName = cc.FullName
	c.Email = cc.Email

	c.CreatedAt = time.Unix(cc.CreatedAt, 0)
	c.UpdatedAt = time.Unix(cc.UpdatedAt, 0)

	return c
}
