package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

var (
	ErrAttendanceQrNotFound = errors.New("attendance QR not found")
)

type AttendanceRepository interface {
	CreateStudentAttendance(ctx context.Context, attendance *models.StudentAttendance) error
	GetStudentAttendancesByStudent(ctx context.Context, studentId uint) ([]models.StudentAttendance, error)
	GetStudentAttendanceHistory(ctx context.Context, studentId uint, from *time.Time, to *time.Time) ([]models.StudentAttendance, error)
	CountAttendanceQrBySectionUntilNow(ctx context.Context, sectionId uint) (int64, error)
	GetAttendanceQrWithSectionAndRoom(ctx context.Context, qrId uint) (*models.AttendanceQr, error)
	GetStudentAttendanceByStudentAndQr(ctx context.Context, studentId uint, qrId uint) (*models.StudentAttendance, error)

	CreateAttendanceQr(ctx context.Context, attendanceQrModelForCreate *AttendanceQrModelForCreate) (*models.AttendanceQr, error)
	DeactivateActiveQrsForSection(ctx context.Context, sectionId uint) error
	GetAttendanceByIDWithSection(ctx context.Context, attendanceId uint) (*models.StudentAttendance, error)
	GetStudentAttendanceBySection(ctx context.Context, sectionId uint, from *time.Time, to *time.Time, studentId *uint) ([]StudentAttendanceWithUser, error)
	UpdateStudentAttendance(ctx context.Context, attendance *models.StudentAttendance) error
	CountStudentAttendanceBySection(ctx context.Context, studentId uint, sectionId uint) (int64, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

type AttendanceQrModelForCreate struct {
	SectionId   uint
	GeneratedAt time.Time
	ExpiresAt   time.Time
}

type StudentAttendanceWithUser = models.StudentAttendance

func (attendanceRepo *attendanceRepository) CreateStudentAttendance(ctx context.Context, attendance *models.StudentAttendance) error {
	if err := attendanceRepo.db.WithContext(ctx).Create(attendance).Error; err != nil {
		return fmt.Errorf("can't create the student attendance: %w", err)
	}
	return nil
}

func (attendanceRepo *attendanceRepository) GetStudentAttendancesByStudent(ctx context.Context, studentId uint) ([]models.StudentAttendance, error) {
	var studentAttendance []models.StudentAttendance
	err := attendanceRepo.db.WithContext(ctx).
		Preload("AttendanceQr").
		Preload("AttendanceQr.CourseSection").
		Preload("AttendanceQr.CourseSection.Course").
		Where("student_id = ?", studentId).
		Order("attended_at DESC").
		Find(&studentAttendance).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch attendances of a student: %w", err)
	}
	return studentAttendance, nil
}

func (attendanceRepo *attendanceRepository) GetStudentAttendanceHistory(ctx context.Context, studentId uint, from *time.Time, to *time.Time) ([]models.StudentAttendance, error) {
	query := attendanceRepo.db.WithContext(ctx).
		Preload("AttendanceQr").
		Preload("AttendanceQr.CourseSection").
		Preload("AttendanceQr.CourseSection.Course").
		Where("student_id = ?", studentId)

	if from != nil {
		query = query.Where("attended_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("attended_at <= ?", *to)
	}

	var studentAttendanceRecords []models.StudentAttendance
	if err := query.Order("attended_at DESC").Find(&studentAttendanceRecords).Error; err != nil {
		return nil, fmt.Errorf("can't fetch attendance history of a student: %w", err)
	}
	return studentAttendanceRecords, nil
}

func (attendanceRepo *attendanceRepository) CountAttendanceQrBySectionUntilNow(ctx context.Context, sectionId uint) (int64, error) {
	var count int64
	now := time.Now().UTC()
	err := attendanceRepo.db.WithContext(ctx).
		Model(&models.AttendanceQr{}).
		Where("section_id = ? AND generated_at <= ?", sectionId, now).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("can't count attendance QR for the section: %w", err)
	}
	return count, nil
}

func (attendanceRepo *attendanceRepository) GetAttendanceQrWithSectionAndRoom(ctx context.Context, qrId uint) (*models.AttendanceQr, error) {
	var attendanceQr models.AttendanceQr
	err := attendanceRepo.db.WithContext(ctx).
		Preload("CourseSection").
		Preload("CourseSection.Room").
		Preload("CourseSection.Semester").
		Preload("CourseSection.SectionSchedules").
		First(&attendanceQr, qrId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttendanceQrNotFound
		}
		return nil, fmt.Errorf("can't find attendance the QR record: %w", err)
	}
	return &attendanceQr, nil
}

func (attendanceRepo *attendanceRepository) GetStudentAttendanceByStudentAndQr(ctx context.Context, studentId uint, qrId uint) (*models.StudentAttendance, error) {
	var attendance models.StudentAttendance
	err := attendanceRepo.db.WithContext(ctx).
		Where("student_id = ? AND qr_id = ?", studentId, qrId).
		First(&attendance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't fetch student attendance for the qr record: %w", err)
	}
	return &attendance, nil
}

func (attendanceRepo *attendanceRepository) CreateAttendanceQr(ctx context.Context, attendanceQrForCreate *AttendanceQrModelForCreate) (*models.AttendanceQr, error) {
	qr := &models.AttendanceQr{
		SectionId:   attendanceQrForCreate.SectionId,
		QrLink:      fmt.Sprintf("SECTION-%d-%d", attendanceQrForCreate.SectionId, attendanceQrForCreate.GeneratedAt.Unix()),
		GeneratedAt: attendanceQrForCreate.GeneratedAt,
		ExpiresAt:   attendanceQrForCreate.ExpiresAt,
		IsActive:    true,
	}
	if err := attendanceRepo.db.WithContext(ctx).Create(qr).Error; err != nil {
		return nil, fmt.Errorf("can't create the attendance QR: %w", err)
	}
	return qr, nil
}

func (attendanceRepo *attendanceRepository) DeactivateActiveQrsForSection(ctx context.Context, sectionId uint) error {
	if err := attendanceRepo.db.WithContext(ctx).
		Model(&models.AttendanceQr{}).
		Where("section_id = ? AND is_active = ?", sectionId, true).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("can't deactivate active the QRs: %w", err)
	}
	return nil
}

func (attendanceRepo *attendanceRepository) GetAttendanceByIDWithSection(ctx context.Context, attendanceId uint) (*models.StudentAttendance, error) {
	var studentAttendance models.StudentAttendance
	err := attendanceRepo.db.WithContext(ctx).
		Preload("User").
		Preload("AttendanceQr").
		Preload("AttendanceQr.CourseSection").
		First(&studentAttendance, attendanceId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("attendance not found: %w", err)
		}
		return nil, fmt.Errorf("can't fetch attendance: %w", err)
	}
	return &studentAttendance, nil
}

func (attendanceRepo *attendanceRepository) GetStudentAttendanceBySection(ctx context.Context, sectionId uint, from *time.Time, to *time.Time, studentId *uint) ([]StudentAttendanceWithUser, error) {
	query := attendanceRepo.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Joins("JOIN attendance_qrs ON attendance_qrs.id = student_attendances.qr_id").
		Where("attendance_qrs.section_id = ?", sectionId).
		Preload("User").
		Preload("AttendanceQr").
		Preload("AttendanceQr.CourseSection").
		Preload("AttendanceQr.CourseSection.Course")

	if from != nil {
		query = query.Where("student_attendances.attended_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("student_attendances.attended_at <= ?", *to)
	}
	if studentId != nil {
		query = query.Where("student_attendances.student_id = ?", *studentId)
	}

	var studentAttendance []models.StudentAttendance
	if err := query.Order("student_attendances.attended_at DESC").Find(&studentAttendance).Error; err != nil {
		return nil, fmt.Errorf("can't fetch attendance for the section: %w", err)
	}
	return studentAttendance, nil
}

func (attendanceRepo *attendanceRepository) UpdateStudentAttendance(ctx context.Context, attendance *models.StudentAttendance) error {
	if err := attendanceRepo.db.WithContext(ctx).Save(attendance).Error; err != nil {
		return fmt.Errorf("can't update the student attendance: %w", err)
	}
	return nil
}

func (attendanceRepo *attendanceRepository) CountStudentAttendanceBySection(ctx context.Context, studentId uint, sectionId uint) (int64, error) {
	var count int64
	err := attendanceRepo.db.WithContext(ctx).
		Model(&models.StudentAttendance{}).
		Joins("JOIN attendance_qrs ON attendance_qrs.id = student_attendances.qr_id").
		Where("student_attendances.student_id = ? AND attendance_qrs.section_id = ? AND student_attendances.attended = ?", studentId, sectionId, true).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("can't count the student attendance: %w", err)
	}
	return count, nil
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}
