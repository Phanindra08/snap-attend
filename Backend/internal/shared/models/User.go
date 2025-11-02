package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName   string      `gorm:"type:varchar(100);not null;index:idx_user_firstname"`
	LastName    string      `gorm:"type:varchar(100);not null"`
	Email       string      `gorm:"type:varchar(100);uniqueIndex:idx_user_email;not null"`
	Password    string      `gorm:"type:varchar(255);not null"`
	ProfileId   uint        `gorm:"not null;index:idx_user_profile"`
	UserProfile UserProfile `gorm:"foreignKey:ProfileId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AddressId   uint        `gorm:"not null;index:idx_user_address"`
	Address     Address     `gorm:"foreignKey:AddressId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
