package dto

type UpdateUserDTO struct {
	FirstName   string
	LastName    string
	OldPassword string
	NewPassword string

	Address1 string
	Address2 string
	City     string
	State    string
	ZipCode  string
	Country  string
}
