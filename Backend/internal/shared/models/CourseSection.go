package models

import "gorm.io/gorm"

type CourseSection struct {
	gorm.Model
	CourseId         uint              `gorm:"not null"`
	Courses          Course            `gorm:"foreignKey:courseId"`
	RoomId           uint              `gorm:"not null"`
	Room             Room              `gorm:"foreignKey:roomId"`
	ProfessorId      uint              `gorm:"not null"`
	ProfessorDetail  User              `gorm:"foreignKey:professorId"`
	SemesterId       uint              `gorm:"not null"`
	Semester         Semester          `gorm:"foreignKey:semesterId"`
	TotalSeats       uint              `gorm:"not null"`
	SectionSchedules []SectionSchedule `gorm:"foreignKey:sectionId"`
	AttendanceQr     []AttendanceQr    `gorm:"foreignKey:sectionId"`
}
