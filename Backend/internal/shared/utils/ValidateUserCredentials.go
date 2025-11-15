package utils

import (
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("date_not_in_past", validateDateNotInPast)
}

// ValidateStructTypes - Helps to validate any struct using the validate struct of Validator package
func ValidateStructTypes(structData interface{}) error {
	return validate.Struct(structData)
}

// validatePassword custom validation for password
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasUpperCaseLetter  = false
		hasLowerCaseLetter  = false
		hasNumber           = false
		hasSpecialCharacter = false
	)
	for _, letter := range password {
		switch {
		case unicode.IsUpper(letter):
			hasUpperCaseLetter = true
		case unicode.IsLower(letter):
			hasLowerCaseLetter = true
		case unicode.IsNumber(letter):
			hasNumber = true
		case unicode.IsPunct(letter) || unicode.IsSymbol(letter):
			hasSpecialCharacter = true
		}
	}
	return len(password) >= 8 && hasUpperCaseLetter && hasLowerCaseLetter && hasNumber && hasSpecialCharacter
}

// validateDateNotInPast - The function helps to validate that the date is not in the past
func validateDateNotInPast(fl validator.FieldLevel) bool {
	var date time.Time
	if datePtr, ok := fl.Field().Interface().(*time.Time); ok {
		if datePtr == nil {
			return true // As nil pointer is considered valid
		}
		date = *datePtr
	} else if dateVal, ok := fl.Field().Interface().(time.Time); ok {
		date = dateVal
	} else {
		return false
	}
	startOfToday := getStartOfToday() // To get the start of today's date
	return !date.Before(startOfToday)
}
