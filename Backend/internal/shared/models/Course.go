package models

type Course struct {
	CourseId   uint   `gorm:"primaryKey; autoIncrement"`
	CourseName string `gorm:"unique; not null"`
}
