package models

type Account struct {
	User          User   `json:"user"`
	Notes         []Note `json:"notes"`
	NumberOfNotes int    `json:"numberOfNotes"`
}
