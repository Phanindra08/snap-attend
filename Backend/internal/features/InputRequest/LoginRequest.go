package InputRequest

import "github.com/phanindra08/snap-attend/internal/shared/models"

type LoginRequest struct {
	Email    string           `json:"email" binding:"required,email"`
	Password string           `json:"password" binding:"required,min=8"`
	Profile  models.UserRoles `json:"profile" binding:"required"`
}
