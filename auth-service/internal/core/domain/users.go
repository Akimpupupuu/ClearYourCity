package core_domain

import (
	"time"
)

var (
	UninitializedID      = -1
	UninitializedVersion = -1
)

type User struct {
	ID        int
	Version   int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

func NewUser(
	id int,
	version int,
	fullName string,
	email string,
	password string,
	createdAt time.Time,
) User {
	return User{
		ID:        id,
		Version:   version,
		FullName:  fullName,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
	}
}

func NewUserUninitialized(fullName string, email string, password string) User {
	return User{
		ID:        UninitializedID,
		Version:   UninitializedVersion,
		FullName:  fullName,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

func (u *User) Validate() error {
	return nil
}
