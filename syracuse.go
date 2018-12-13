package syracuse

import "time"

// Citizen ...
type Citizen struct {
	ID       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Fullname string `json:"fullname" db:"fullname"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// Citizens ...
type Citizens interface {
	Get(string) (*Citizen, error)
	Select() ([]*Citizen, error)
	Create(*Citizen) error
	Update(*Citizen) error
	Delete(*Citizen) error
}
