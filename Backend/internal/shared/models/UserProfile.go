package models

type userRoles string

const (
	student   userRoles = "Student"
	professor userRoles = "Professor"
	admin     userRoles = "Admin"
)

type UserProfile struct {
	ProfileId uint      `gorm:"primaryKey; autoIncrement"`
	Role      userRoles `gorm:"unique; not null"`
	User      []User
}
