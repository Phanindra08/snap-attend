package models

import "gorm.io/datatypes"

type StudentAttendance struct {
	attendanceId   uint           `gorm:"primaryKey"`
	studentId      uint           `gorm:"not null"`
	courseSection  CourseSection  `gorm:"foreignKey:sectionId"`
	qrLink         string         `gorm:"not null"`
	generatedDate  datatypes.Date `gorm:"not null"`
	generatedTime  datatypes.Time `gorm:"not null"`
	expirationTime datatypes.Time `gorm:"not null"`
}
