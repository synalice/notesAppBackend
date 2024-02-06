package models

import "time"

type Note struct {
	ID          int       `json:"id"`
	AuthorID    int       `json:"authorID"`
	DateCreated time.Time `json:"dateCreated"`
	Symbols     int       `json:"symbols"`
	Contents    string    `json:"contents"`
	Title       string    `json:"title"`
}
