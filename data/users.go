package data

import (
	"errors"
	"time"

	"database/sql"

	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidEmailOrPassword = errors.New("tokens: Invalid email or password")

// NewUser is the type received from mobile apps before being saved into the db
type NewUser struct {
	Name     string
	Email    string
	Password string
}

// User contains info on our user
type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password []byte
	Created  time.Time
}

// UserStore is our interface defining methods for concrete service
type UserStore interface {
	GetUser(id string) (User, error)
}

// UserService is our psql implementation of UserStore
type UserService struct {
	db *sql.DB
}

// InitUserService initializes a new UserService
func InitUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// NewUser inserts a new user into the database
func (us *UserService) NewUser(nu NewUser) (User, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("users - NewUser: %v", err)
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return User{}, fmt.Errorf("users - NewUser: failed to generate id: %v", err)
	}

	u := User{
		ID:       id,
		Name:     nu.Name,
		Email:    nu.Email,
		Password: hp,
		Created:  time.Now().Round(time.Second),
	}

	_, err = us.db.Exec(
		"Insert into users (user_id, email, name, password, date_created) VALUES ($1, $2, $3, $4, $5)",
		u.ID,
		u.Email,
		u.Name,
		u.Password,
		u.Created,
	)

	if err != nil {
		return User{}, fmt.Errorf("users - NewUser: %v", err)
	}

	return u, nil
}

// NewAnonUser creates a new anon user without an details other than an id
func (us *UserService) NewAnonUser() (User, error) {
	u := User{
		ID:      uuid.New(),
		Created: time.Now().Round(time.Second),
	}

	a, err := us.db.Exec(
		"Insert into users (user_id, date_created) VALUES ($1, $2)",
		u.ID,
		u.Created,
	)
	fmt.Println(a)

	if err != nil {
		return User{}, fmt.Errorf("users - NewAnonUser: %v", err)
	}

	return u, nil
}

// GetUser returns a user from the database
func (us *UserService) GetUser(id string) (User, error) {
	var email sql.NullString
	var name sql.NullString

	row := us.db.QueryRow("SELECT user_id, email, name, password, date_created FROM users WHERE user_id = $1", id)

	u := User{}

	err := row.Scan(&u.ID, &email, &name, &u.Password, &u.Created)

	if email.Valid {
		u.Email = email.String
	} else {
		u.Email = ""
	}

	if name.Valid {
		u.Name = name.String
	} else {
		u.Name = ""
	}

	if err != nil {
		return User{}, err
	}

	// Convert time to local
	u.Created = u.Created.In(time.Local)

	return u, nil
}

// Authenticate returns a user if the given email and password match a user in the db
func (us *UserService) Authenticate(e, p string) (User, error) {
	row := us.db.QueryRow("SELECT user_id, email, name, password, date_created FROM users WHERE email = $1", e)

	u := User{}

	err := row.Scan(&u.ID, &u.Email, &u.Email, &u.Password, &u.Created)

	if err != nil {

		if err == sql.ErrNoRows {
			return User{}, ErrInvalidEmailOrPassword
		}

		return User{}, fmt.Errorf("users - Authenticate: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(p))

	if err != nil {
		return User{}, ErrInvalidEmailOrPassword
	}

	return u, nil
}
