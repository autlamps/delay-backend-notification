package data

import (
	"database/sql"
	"time"
)

// User contains info on our user
type User struct {
	ID       string
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

// UserStore is our interface defining methods for concrete service
type UserStore interface {
}

// UserService is our psql implementation of UserStore
type UserService struct {
	db *sql.DB
}
