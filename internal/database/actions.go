package database

import (
	"database/sql"
	"errors"
	"fmt"
	"notesAppBackend/internal/models"
)

// TODO: Read more about SQL queries (https://sqlbolt.com/).

// GetUserByEmail takes email and return a user from the database.
func (db *Database) GetUserByEmail(email string) (models.User, error) {
	stmt, err := db.connection.Prepare(`SELECT * FROM "user" WHERE "user".email = $1`)
	if err != nil {
		return models.User{}, fmt.Errorf("GetUserByEmail: %w", err)
	}

	var user models.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.HashedPassword, &user.DateCreated, &user.Nickname, &user.PFPLink)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("GetUserByEmail: %w", err)
	}

	return user, nil
}

// CreateNewUser is responsible for saving new user to the database.
func (db *Database) CreateNewUser(user models.User) error {
	if user == (models.User{}) {
		return fmt.Errorf("CreateNewUser: user struct empty")
	}

	// Start the transaction.
	tx, err := db.connection.Begin()
	if err != nil {
		return fmt.Errorf("CreateNewUser: %w", err)
	}

	//Defer a rollback in case anything fails.
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	stmt, err := tx.Prepare(`INSERT INTO "user" (email, hashed_password, nickname, date_created) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return fmt.Errorf("CreateNewUser: %w", err)
	}

	_, err = stmt.Exec(user.Email, user.HashedPassword, user.Nickname, user.DateCreated)
	if err != nil {
		return fmt.Errorf("CreateNewUser: %w", err)
	}

	// Commit the transaction.
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("CreateNewUser: %w", err)
	}

	return nil
}
