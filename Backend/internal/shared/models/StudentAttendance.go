package models

import (
	"time"
)

type StudentAttendance struct {
	AttendanceId    uint         `gorm:"primaryKey;autoIncrement"`
	StudentId       uint         `gorm:"not null;uniqueIndex:idx_attendance_unique;index:idx_attendance_student"`
	User            User         `gorm:"foreignKey:StudentId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	QrId            uint         `gorm:"not null;uniqueIndex:idx_attendance_unique;index:idx_attendance_qr"`
	AttendanceQr    AttendanceQr `gorm:"foreignKey:QrId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AttendedAt      time.Time    `gorm:"not null"`
	Attended        bool         `gorm:"not null;default:false"`
	StudentQuestion string       `gorm:"type:text"`
	StudentAnswer   string       `gorm:"type:text"`
}
