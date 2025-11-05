package models

import "gorm.io/gorm"

type CourseSection struct {
	gorm.Model
	SectionNumber    string            `gorm:"type:varchar(10);not null;uniqueIndex:idx_section_unique"`
	CourseId         uint              `gorm:"not null;uniqueIndex:idx_section_unique;index:idx_section_course"`
	Course           Course            `gorm:"foreignKey:CourseId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	RoomId           uint              `gorm:"not null;index:idx_section_room"`
	Room             Room              `gorm:"foreignKey:RoomId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ProfessorId      uint              `gorm:"not null;index:idx_section_professor"`
	ProfessorDetail  User              `gorm:"foreignKey:ProfessorId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SemesterId       uint              `gorm:"not null;uniqueIndex:idx_section_unique;index:idx_section_semester"`
	Semester         Semester          `gorm:"foreignKey:SemesterId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalSeats       uint              `gorm:"not null"`
	EnrolledSeats    uint              `gorm:"default:0"`
	SectionSchedules []SectionSchedule `gorm:"foreignKey:SectionId"`
	AttendanceQrs    []AttendanceQr    `gorm:"foreignKey:SectionId"`
}
