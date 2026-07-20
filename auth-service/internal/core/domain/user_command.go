package core_domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/errors"
)

var (
	RegularExpression = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
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

type PatchPasswordCommand struct {
	OldPassword string
	NewPassword string
}

type PatchUserCommand struct {
	FullName *string
	Email    *string
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

func NewPatchPasswordCommand(oldPassword string, newPassword string) PatchPasswordCommand {
	return PatchPasswordCommand{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
}

func NewPatchUserCommand(fullName *string, email *string) PatchUserCommand {
	return PatchUserCommand{
		FullName: fullName,
		Email:    email,
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
		return fmt.Errorf("invalid 'password' length: %d: %w", passwordLength, core_errors.ErrInvalidArgument)
	}

	return nil
}

func (c *LoginCommand) Validate() error {
	if c.Password == "" {
		return fmt.Errorf("'password' is required: %w", core_errors.ErrInvalidArgument)
	}

	emailLength := len([]rune(c.Email))
	if emailLength < 5 || emailLength > 100 {
		return fmt.Errorf("invalid 'email' length: %d: %w", emailLength, core_errors.ErrInvalidArgument)
	}

	return nil
}

func (c *PatchPasswordCommand) Validate() error {
	if c.OldPassword == "" {
		return fmt.Errorf("'oldPassword' is required: %w", core_errors.ErrInvalidArgument)
	}

	newPasswordLength := len([]rune(c.NewPassword))
	if newPasswordLength < 8 {
		return fmt.Errorf("invalid 'newPassword' length: %d: %w", newPasswordLength, core_errors.ErrInvalidArgument)
	}

	return nil
}

func (c *PatchUserCommand) Validate() error {
	if c.FullName != nil {
		fullNameLength := len([]rune(*c.FullName))

		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("invalid 'fullName' length: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
		}
	}

	if c.Email != nil {
		emailLength := len([]rune(*c.Email))

		if emailLength < 5 || emailLength > 100 {
			return fmt.Errorf("invalid 'email' length: %d: %w", emailLength, core_errors.ErrInvalidArgument)
		}

		if !RegularExpression.MatchString(*c.Email) {
			return fmt.Errorf("invalid 'email' format: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
