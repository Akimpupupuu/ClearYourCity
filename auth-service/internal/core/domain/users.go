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

type RegisterCommand struct {
	FullName string
	Email    string
	Password string
}

type LoginCommand struct {
	Email    string
	Password string
}

func NewUser(
	id int,
	version int,
	fullName string,
	email string,
	hashedPassword string,
	createdAt time.Time,
) User {
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

func NewRegisterCommand(fullName string, email string, password string) RegisterCommand {
	return RegisterCommand{
		FullName: fullName,
		Email:    email,
		Password: password,
	}
}

func NewLoginCommand(email string, password string) LoginCommand {
	return LoginCommand{
		Email:    email,
		Password: password,
	}
}

func (c *RegisterCommand) Validate() error {
	fullNameLength := len([]rune(c.FullName))
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid 'fullName' length: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	emailLength := len(c.Email)
	if emailLength < 5 || emailLength > 100 {
		return fmt.Errorf("invalid 'email' length: %d: %w", emailLength, core_errors.ErrInvalidArgument)
	}

	if !RegularExpression.MatchString(c.Email) {
		return fmt.Errorf("invalid 'email' format: %w", core_errors.ErrInvalidArgument)
	}

	passwordLength := len([]rune(c.Password))
	if passwordLength < 8 {
		return fmt.Errorf("invalid 'password' length: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
