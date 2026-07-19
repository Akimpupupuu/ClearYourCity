package core_domain

import (
	"fmt"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

type RegisterCommand struct {
	FullName string
	Email    string
	Password string
}

type LoginCommand struct {
	Email    string
	Password string
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
