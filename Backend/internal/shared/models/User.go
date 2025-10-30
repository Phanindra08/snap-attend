package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string `gorm:"not null"`
	LastName    string `gorm:"not null"`
	Email       string `gorm:"unique; not null"`
	Password    string `gorm:"not null"`
	ProfileId   uint
	UserProfile UserProfile `gorm:"foreignKey:ProfileId"`
	AddressId   uint
	Address     Address `gorm:"foreignKey:AddressId; constraint: OnUpdate:CASCADE, OnDelete:CASCADE"`
}
