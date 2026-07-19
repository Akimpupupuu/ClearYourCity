package core_domain

import (
	"fmt"
	"regexp"
	"time"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

var (
	UninitializedID      = -1
	UninitializedVersion = -1
	RegularExpression    = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

type User struct {
	ID             int
	Version        int
	FullName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

func NewUser(id int, version int, fullName string, email string, hashedPassword string, createdAt time.Time) User {
	return User{
		ID:             id,
		Version:        version,
		FullName:       fullName,
		Email:          email,
		HashedPassword: hashedPassword,
		CreatedAt:      createdAt,
	}
}

func NewUserUninitialized(fullName string, email string, hashedPassword string) User {
	return User{
		ID:             UninitializedID,
		Version:        UninitializedVersion,
		FullName:       fullName,
		Email:          email,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}
}

func (u *User) ApplyPatch(fullName *string, email *string) error {
	tmp := *u

	if fullName != nil {
		tmp.FullName = *fullName
	}

	if email != nil {
		tmp.Email = *email
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate user: %w", err)
	}

	*u = tmp
	return nil
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid 'fullName' length: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	emailLength := len(u.Email)
	if emailLength < 5 || emailLength > 100 {
		return fmt.Errorf("invalid 'email' length: %d: %w", emailLength, core_errors.ErrInvalidArgument)
	}

	if !RegularExpression.MatchString(u.Email) {
		return fmt.Errorf("invalid 'email' format: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
