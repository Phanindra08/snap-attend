package models

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
)

const (
	SALT_LENGTH = 16
	TIME_COST   = 3
	MEMORY_COST = 64 * 1024 // Using 64 MB
	THREADS     = 4
	KEY_LENGTH  = 32
)

type User struct {
	gorm.Model
	FirstName          string              `gorm:"type:varchar(100);not null;index:idx_user_name"`
	LastName           string              `gorm:"type:varchar(100);not null;index:idx_user_name"`
	Email              string              `gorm:"type:varchar(100);uniqueIndex:idx_user_email;not null"`
	Password           string              `gorm:"type:varchar(255);not null"`
	ProfileId          uint                `gorm:"not null;index:idx_user_profile"`
	UserProfile        UserProfile         `gorm:"foreignKey:ProfileId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AddressId          uint                `gorm:"not null;index:idx_user_address"`
	Address            Address             `gorm:"foreignKey:AddressId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CourseSections     []CourseSection     `gorm:"foreignKey:ProfessorId"`
	StudentAttendances []StudentAttendance `gorm:"foreignKey:StudentId"`
}

// PasswordHashing - Helpful to hash the user's password using Argon2 package
func (user *User) PasswordHashing() error {
	if len(user.Password) == 0 {
		return errors.New("password cannot be empty")
	}

	// Checking the strength of the Password
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	// Generating a random salt
	salt := make([]byte, SALT_LENGTH)
	_, err := rand.Read(salt)
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}

	// Hashing the password with Argon2
	hashedPassword := argon2.IDKey([]byte(user.Password), salt, TIME_COST, MEMORY_COST, THREADS, KEY_LENGTH)

	// Store the salt and hashed password as a single string
	user.Password = fmt.Sprintf("%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashedPassword))
	return nil
}

// ComparePassword - compares the provided password with the stored hash
func (user *User) ComparePassword(password string) bool {
	if len(user.Password) == 0 || len(password) == 0 {
		return false
	}

	// Splitting the stored password into salt and hash
	parts := strings.Split(user.Password, "$")
	if len(parts) != 2 {
		return false
	}

	// Decoding the salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// Hash the provided password with the same salt
	computedHash := argon2.IDKey([]byte(password), salt, TIME_COST, MEMORY_COST, THREADS, KEY_LENGTH)

	// Compare the computed hash with the stored hash
	return subtle.ConstantTimeCompare(storedHash, computedHash) == 1
}
