package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.StudentSectionEnrollment) error
	GetEnrollmentByStudentAndSection(ctx context.Context, studentId uint, sectionId uint) (*models.StudentSectionEnrollment, error)
	GetEnrollmentsByStudentAndSemester(ctx context.Context, studentId uint, semesterId uint) ([]models.StudentSectionEnrollment, error)
	CountEnrollmentsByStudentAndSemester(ctx context.Context, studentId uint, semesterId uint) (int64, error)
	GetEnrollmentsBySection(ctx context.Context, sectionId uint) ([]models.StudentSectionEnrollment, error)
}

type enrollmentRepository struct {
	db *gorm.DB
}

func (enrollRepo *enrollmentRepository) CreateEnrollment(ctx context.Context, enrollment *models.StudentSectionEnrollment) error {
	if err := enrollRepo.db.WithContext(ctx).Create(enrollment).Error; err != nil {
		return fmt.Errorf("can't create an enrollment: %w", err)
	}
	return nil
}

func (enrollRepo *enrollmentRepository) GetEnrollmentByStudentAndSection(ctx context.Context, studentId uint, sectionId uint) (*models.StudentSectionEnrollment, error) {
	var enrollment models.StudentSectionEnrollment
	err := enrollRepo.db.WithContext(ctx).
		Where("student_id = ? AND section_id = ?", studentId, sectionId).
		First(&enrollment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("can't fetch an enrollment: %w", err)
	}
	return &enrollment, nil
}

func (enrollRepo *enrollmentRepository) GetEnrollmentsByStudentAndSemester(ctx context.Context, studentId uint, semesterId uint) ([]models.StudentSectionEnrollment, error) {
	var enrollments []models.StudentSectionEnrollment
	err := enrollRepo.db.WithContext(ctx).
		Preload("CourseSection").
		Preload("CourseSection.SectionSchedules").
		Preload("CourseSection.Course").
		Preload("CourseSection.Semester").
		Joins("JOIN course_sections ON course_sections.id = student_section_enrollments.section_id").
		Where("student_section_enrollments.student_id = ? AND course_sections.semester_id = ?", studentId, semesterId).
		Find(&enrollments).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch enrollments: %w", err)
	}
	return enrollments, nil
}

func (enrollRepo *enrollmentRepository) CountEnrollmentsByStudentAndSemester(ctx context.Context, studentId uint, semesterId uint) (int64, error) {
	var count int64
	err := enrollRepo.db.WithContext(ctx).
		Table("student_section_enrollments").
		Joins("JOIN course_sections ON course_sections.id = student_section_enrollments.section_id").
		Where("student_section_enrollments.student_id = ? AND course_sections.semester_id = ?", studentId, semesterId).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("can't count enrollments: %w", err)
	}
	return count, nil
}

func (enrollRepo *enrollmentRepository) GetEnrollmentsBySection(ctx context.Context, sectionId uint) ([]models.StudentSectionEnrollment, error) {
	var enrollments []models.StudentSectionEnrollment
	err := enrollRepo.db.WithContext(ctx).
		Preload("Student").
		Where("section_id = ?", sectionId).
		Find(&enrollments).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch enrollments by section: %w", err)
	}
	return enrollments, nil
}

func NewEnrollmentRepository(db *gorm.DB) EnrollmentRepository {
	return &enrollmentRepository{db: db}
}
