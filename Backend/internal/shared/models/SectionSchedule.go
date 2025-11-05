package models

import "gorm.io/datatypes"

type DayOfWeek string

const (
	Monday    DayOfWeek = "Monday"
	Tuesday   DayOfWeek = "Tuesday"
	Wednesday DayOfWeek = "Wednesday"
	Thursday  DayOfWeek = "Thursday"
	Friday    DayOfWeek = "Friday"
	Saturday  DayOfWeek = "Saturday"
	Sunday    DayOfWeek = "Sunday"
)

type SectionSchedule struct {
	ScheduleId     uint           `gorm:"primaryKey;autoIncrement"`
	SectionId      uint           `gorm:"not null;uniqueIndex:idx_schedule_unique;index:idx_schedule_section"`
	CourseSection  CourseSection  `gorm:"foreignKey:SectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DaysOfTheClass DayOfWeek      `gorm:"type:varchar(15);not null;uniqueIndex:idx_schedule_unique"`
	StartTime      datatypes.Time `gorm:"not null;uniqueIndex:idx_schedule_unique;index:idx_schedule_start_time"`
	EndTime        datatypes.Time `gorm:"not null"`
}
