package models

import (
	"time"
)

type AttendanceQr struct {
	QrId               uint                `gorm:"primaryKey;autoIncrement"`
	SectionId          uint                `gorm:"not null;index:idx_qr_section"`
	CourseSection      CourseSection       `gorm:"foreignKey:SectionId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	QrLink             string              `gorm:"type:varchar(255);uniqueIndex;not null"`
	GeneratedAt        time.Time           `gorm:"not null"`
	ExpiresAt          time.Time           `gorm:"not null"`
	IsActive           bool                `gorm:"default:true;index:idx_qr_is_active"`
	StudentAttendances []StudentAttendance `gorm:"foreignKey:QrId"`
}
