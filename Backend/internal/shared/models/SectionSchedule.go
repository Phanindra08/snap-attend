package models

import "gorm.io/datatypes"

type daysOfTheWeek string

const (
	monday    daysOfTheWeek = "Monday"
	tuesday   daysOfTheWeek = "Tuesday"
	wednesday daysOfTheWeek = "Wednesday"
	thursday  daysOfTheWeek = "Thursday"
	friday    daysOfTheWeek = "Friday"
	saturday  daysOfTheWeek = "Saturday"
	sunday    daysOfTheWeek = "Sunday"
)

type SectionSchedule struct {
	scheduleId     uint           `gorm:"primaryKey; autoIncrement"`
	sectionId      uint           `gorm:"not null"`
	courseSection  CourseSection  `gorm:"foreignKey:sectionId"`
	daysOfTheClass daysOfTheWeek  `gorm:"not null"`
	startTime      datatypes.Time `gorm:"not null"`
	endTime        datatypes.Time `gorm:"not null"`
}
