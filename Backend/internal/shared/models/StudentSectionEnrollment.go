package models

import "time"

type StudentSectionEnrollment struct {
	ID            uint          `gorm:"primaryKey;autoIncrement"`
	StudentId     uint          `gorm:"not null;uniqueIndex:idx_enrollment_student_section;index:idx_enrollment_student"`
	Student       User          `gorm:"foreignKey:StudentId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SectionId     uint          `gorm:"not null;uniqueIndex:idx_enrollment_student_section;index:idx_enrollment_section"`
	CourseSection CourseSection `gorm:"foreignKey:SectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EnrolledAt    time.Time     `gorm:"not null"`
}
