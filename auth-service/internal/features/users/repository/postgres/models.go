package users_postrges_repository

import "time"

type UserModel struct {
	ID             int
	Version        int
	FullName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}
