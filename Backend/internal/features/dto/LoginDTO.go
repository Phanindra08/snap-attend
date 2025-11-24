package dto

import "github.com/phanindra08/snap-attend/internal/shared/models"

type LoginDTO struct {
	Email    string
	Password string
	Profile  models.UserRoles
}
