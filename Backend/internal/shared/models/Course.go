package models

type Course struct {
	CourseId       uint            `gorm:"primaryKey;autoIncrement"`
	CourseName     string          `gorm:"type:varchar(100);uniqueIndex:idx_course_name;not null"`
	CourseSections []CourseSection `gorm:"foreignKey:CourseId"`
}
