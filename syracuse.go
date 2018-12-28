package syracuse

import "time"

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
	Get(*CitizensQuery) (*Citizen, error)
	Select() ([]*Citizen, error)
	Create(*Citizen) error
	Update(*Citizen) error
	Delete(*Citizen) error
}

// CitizensStore crud methods to make queries over some database. Any struct that implementes this methods could use Citizen.
type CitizensStore interface {
	Get(*CitizensQuery) (*Citizen, error)
	Select() ([]*Citizen, error)
	Create(*Citizen) error
	Update(*Citizen) error
	Delete(*Citizen) error
}

// CitizensQuery represents queries that helps to query users.
type CitizensQuery struct {
	ID       string
	Email    string
	FullName string
}
