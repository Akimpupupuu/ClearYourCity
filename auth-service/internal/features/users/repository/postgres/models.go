package users_repository

import (
	"time"

	core_domain "github.com/Akimpupupuu/ClearYourCity/auth-service/internal/core/domain"
)

type UserModel struct {
	ID             int
	Version        int
	FullName       string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}

func DomainFromModel(model UserModel) core_domain.User {
	return core_domain.NewUser(
		model.ID,
		model.Version,
		model.FullName,
		model.Email,
		model.HashedPassword,
		model.CreatedAt,
	)
}
