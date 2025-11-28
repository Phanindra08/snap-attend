package InputRequest

type AdminCreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`

	Address1 string `json:"address1" binding:"required"`
	Address2 string `json:"address2"`
	City     string `json:"city" binding:"required"`
	State    string `json:"state" binding:"required"`
	ZipCode  string `json:"zipCode" binding:"required"`
	Country  string `json:"country"`
}

type AdminUpdateUserRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`

	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	City     string `json:"city,omitempty"`
	State    string `json:"state,omitempty"`
	ZipCode  string `json:"zipCode,omitempty"`
	Country  string `json:"country,omitempty"`
}
