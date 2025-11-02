package models

type Address struct {
	AddressId uint   `gorm:"primaryKey;autoIncrement"`
	Address1  string `gorm:"type:varchar(255);not null"`
	Address2  string `gorm:"type:varchar(100)"`
	City      string `gorm:"type:varchar(100);not null"`
	State     string `gorm:"type:varchar(50);not null"`
	ZipCode   string `gorm:"type:varchar(15);not null"`
	Country   string `gorm:"type:varchar(50);not null;default:'USA'"`
	Users     []User `gorm:"foreignKey:AddressId"`
}
