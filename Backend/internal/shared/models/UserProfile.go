package models

type UserRoles string

const (
	Student   UserRoles = "Student"
	Professor UserRoles = "Professor"
	Admin     UserRoles = "Admin"
)

type UserProfile struct {
	ID    uint      `gorm:"primaryKey;autoIncrement"`
	Role  UserRoles `gorm:"type:varchar(15);uniqueIndex;not null"`
	Users []User    `gorm:"foreignKey:ProfileId"`
}
