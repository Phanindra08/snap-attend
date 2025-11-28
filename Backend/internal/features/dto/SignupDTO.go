package dto

import "github.com/phanindra08/snap-attend/internal/shared/models"

type SignupDTO struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Profile   models.UserRoles

	Address1 string
	Address2 string
	City     string
	State    string
	ZipCode  string
	Country  string
}
