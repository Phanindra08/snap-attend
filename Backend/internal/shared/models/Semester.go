package models

import "gorm.io/gorm"
import "gorm.io/datatypes"

type typesOfSemesters string

const (
	fall    typesOfSemesters = "Fall"
	winter  typesOfSemesters = "Winter"
	summer1 typesOfSemesters = "Summer-1"
	summer2 typesOfSemesters = "Summer-2"
)

type Semester struct {
	semesterId   uint             `gorm:"primaryKey"`
	semesterType typesOfSemesters `gorm:"not null"`
	year         uint16           `gorm:"not null"`
	startDate    datatypes.Date   `gorm:"not null"`
	endDate      datatypes.Date   `gorm:"not null"`
	gorm.Model
}
