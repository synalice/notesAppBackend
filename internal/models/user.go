package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	HashedPassword string    `json:"hashedPassword,omitempty"`
	DateCreated    time.Time `json:"dateCreated,omitempty"`
	Nickname       string    `json:"nickname,omitempty"`

	// TODO: `omitempty` does work. There needs to be a custom type for that will include `sql.NullString`
	// 	inside itself so that it's possible to redefine `MarshalJSON` method.
	PFPLink sql.NullString `json:"pfpLink,omitempty"` // Link to a profile picture
}
