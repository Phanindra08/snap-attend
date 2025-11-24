package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SemesterTypes string

const (
	Fall    SemesterTypes = "Fall"
	Winter  SemesterTypes = "Winter"
	Summer1 SemesterTypes = "Summer-1"
	Summer2 SemesterTypes = "Summer-2"
)

type Semester struct {
	gorm.Model
	SemesterType   SemesterTypes   `gorm:"type:varchar(15);not null;uniqueIndex:idx_semester_unique"`
	Year           uint16          `gorm:"not null;uniqueIndex:idx_semester_unique;index:idx_semester_year"`
	StartDate      datatypes.Date  `gorm:"not null"`
	EndDate        datatypes.Date  `gorm:"not null"`
	CourseSections []CourseSection `gorm:"foreignKey:SemesterId"`
}
