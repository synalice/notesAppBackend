package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	connection *sql.DB
}

func New(connectionURL string) (*Database, error) {
	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		log.Println("Couldn't create a database connection")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Couldn't Ping() a database connection")
		return nil, err
	}

	database := Database{
		connection: db,
	}

	return &database, nil
}
