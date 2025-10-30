package models

type Address struct {
	AddressId uint   `gorm:"primaryKey; autoIncrement"`
	Address1  string `gorm:"not null"`
	Address2  string
	City      string `gorm:"not null"`
	State     string `gorm:"not null"`
	ZipCode   string `gorm:"not null"`
	Country   string `gorm:"not null"`
}
