package models

type Room struct {
	RoomId         uint            `gorm:"primaryKey;autoIncrement"`
	RoomNumber     string          `gorm:"type:varchar(15);not null;uniqueIndex:idx_room_unique"`
	BuildingName   string          `gorm:"type:varchar(50);not null;uniqueIndex:idx_room_unique"`
	Capacity       int             `gorm:"not null"`
	Address        string          `gorm:"type:varchar(255);not null"`
	CourseSections []CourseSection `gorm:"foreignKey:RoomId"`
}
