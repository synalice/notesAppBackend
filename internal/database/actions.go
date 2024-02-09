package database

import (
	"database/sql"
	"errors"
	"fmt"
	"notesAppBackend/internal/models"
	"time"
)

// TODO: Read more about SQL queries (https://sqlbolt.com/).

// GetUserByNickname return a user with specified nickname from the database.
func (db *Database) GetUserByNickname(email string) (models.User, error) {
	stmt, err := db.connection.Prepare(`SELECT * FROM "user" WHERE "user".nickname = $1`)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByNickname: %w", err)
	}

	var user models.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.HashedPassword, &user.DateCreated, &user.Nickname, &user.PFPLink)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("GetUserByNickname: %w", err)
	}

	return user, nil
}

// CreateNewUser is responsible for saving new user to the database.
func (db *Database) CreateNewUser(user models.User) (newUserID int, err error) {
	if user == (models.User{}) {
		return 0, fmt.Errorf("CreateNewUser: user struct empty")
	}

	// Start the transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		return 0, fmt.Errorf("CreateNewUser: %w", err)
	}

	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.Prepare(`INSERT INTO "user" (hashed_password, nickname, date_created) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("CreateNewUser: %w", err)
	}

	err = stmt.QueryRow(user.HashedPassword, user.Nickname, user.DateCreated).Scan(&newUserID)
	if err != nil {
		return 0, fmt.Errorf("CreateNewUser: %w", err)
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("CreateNewUser: %w", err)
	}

	return newUserID, nil
}

func (db *Database) CreateNewPost(authorID int, title, contents string) error {
	contentsLen := len(contents)

	if len(title) <= 0 {
		return fmt.Errorf("CreateNewPost: title is empty")
	}

	if contentsLen <= 0 {
		return fmt.Errorf("CreateNewPost: contents are empty")
	}

	// Start the transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		return fmt.Errorf("CreateNewPost: %w", err)
	}

	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.Prepare(`INSERT INTO "note" (author_id, date_created, symbols, contents, title) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return fmt.Errorf("CreateNewPost: %w", err)
	}

	_, err = stmt.Exec(authorID, time.Now(), contentsLen, contents, title)
	if err != nil {
		return fmt.Errorf("CreateNewPost: %w", err)
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CreateNewPost: %w", err)
	}

	return nil
}

func (db *Database) GetAccountData(accountID int) (models.Account, error) {
	// Start the transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	// Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmtUser, err := tx.Prepare(`SELECT id, nickname, date_created FROM "user" WHERE id = $1`)
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	stmtNotes, err := tx.Prepare(`SELECT * FROM "note" WHERE author_id = $1 ORDER BY date_created DESC`)
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	var user models.User
	var notes []models.Note

	err = stmtUser.QueryRow(accountID).Scan(&user.ID, &user.Nickname, &user.DateCreated)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.Account{}, nil
		}
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	rows, err := stmtNotes.Query(accountID)
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Loop through rows.
	for rows.Next() {
		var note models.Note

		err = rows.Scan(&note.ID, &note.AuthorID, &note.DateCreated, &note.Symbols, &note.Contents, &note.Title)
		if err != nil {
			return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
		}

		notes = append(notes, note)
	}

	err = rows.Err()
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return models.Account{}, fmt.Errorf("GetAccountData: %w", err)
	}

	return models.Account{
		User:          user,
		Notes:         notes,
		NumberOfNotes: len(notes),
	}, nil
}

func (db *Database) IsIDPresent(accountID int) (bool, error) {
	stmt, err := db.connection.Prepare(`SELECT id FROM "user" WHERE "user".id = $1`)
	if err != nil {
		return false, fmt.Errorf("IsIDPresent: %w", err)
	}

	var id int
	err = stmt.QueryRow(accountID).Scan(&id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return false, nil
		}
		return false, fmt.Errorf("IsIDPresent: %w", err)
	}

	return true, nil
}

func (db *Database) FindNoteByID(noteID int) (models.Note, error) {
	stmt, err := db.connection.Prepare(`SELECT id, date_created, author_id, contents, symbols, title FROM "note" WHERE "note".id = $1`)
	if err != nil {
		return models.Note{}, fmt.Errorf("FindNoteByID: %w", err)
	}

	var note models.Note

	err = stmt.QueryRow(noteID).Scan(&note.ID, &note.DateCreated, &note.AuthorID, &note.Contents, &note.Symbols, &note.Title)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.Note{}, nil
		}
		return models.Note{}, fmt.Errorf("IsIDPresent: %w", err)
	}

	return note, nil
}

func (db *Database) DeleteNoteByID(noteID int) error {
	stmt, err := db.connection.Prepare(`DELETE FROM "note" WHERE "note".id = $1`)
	if err != nil {
		return fmt.Errorf("DeleteNoteByID: %w", err)
	}

	_, err = stmt.Exec(noteID)
	if err != nil {
		return fmt.Errorf("DeleteNoteByID: %w", err)
	}

	return nil
}
