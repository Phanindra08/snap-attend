package models

type Room struct {
	RoomId       uint   `gorm:"primaryKey; autoIncrement"`
	RoomNumber   string `gorm:"not null"`
	BuildingName string `gorm:"not null"`
	Location     string `gorm:"not null"`
}
