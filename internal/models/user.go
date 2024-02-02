package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int            `json:"id"`
	Email          string         `json:"email"`
	HashedPassword string         `json:"hashedPassword"`
	DateCreated    time.Time      `json:"dateCreated"`
	Nickname       string         `json:"nickname"`
	PFPLink        sql.NullString `json:"pfpLink"` // Link to a profile picture
}
