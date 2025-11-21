package InputRequest

import "github.com/phanindra08/snap-attend/internal/shared/models"

type SignupRequest struct {
	FirstName string           `json:"firstName" binding:"required"`
	LastName  string           `json:"lastName" binding:"required"`
	Email     string           `json:"email" binding:"required,email"`
	Password  string           `json:"password" binding:"required,min=8"`
	Profile   models.UserRoles `json:"profile" binding:"required"`

	Address1 string `json:"address1" binding:"required"`
	Address2 string `json:"address2"`
	City     string `json:"city" binding:"required"`
	State    string `json:"state" binding:"required"`
	ZipCode  string `json:"zipCode" binding:"required"`
	Country  string `json:"country" binding:"required"`
}
