package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

var (
	ErrSectionNotFound = errors.New("section not found")
)

type SectionRepository interface {
	GetSectionByID(ctx context.Context, id uint) (*models.CourseSection, error)
	GetSectionsByProfessorAndSemester(ctx context.Context, professorId uint, semesterId uint) ([]models.CourseSection, error)
	GetSectionsByRoomAndSemester(ctx context.Context, roomId uint, semesterId uint) ([]models.CourseSection, error)
	UpdateSection(ctx context.Context, section *models.CourseSection) error
	GetSectionsByProfessor(ctx context.Context, professorId uint) ([]models.CourseSection, error)
	SearchSectionsByProfessorAndCourseName(ctx context.Context, professorId uint, name string) ([]models.CourseSection, error)
}

type sectionRepository struct {
	db *gorm.DB
}

func (sectionRepo *sectionRepository) GetSectionByID(ctx context.Context, id uint) (*models.CourseSection, error) {
	var section models.CourseSection
	if err := sectionRepo.db.WithContext(ctx).
		Preload("SectionSchedules").
		Preload("Course").
		Preload("Semester").
		First(&section, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSectionNotFound
		}
		return nil, fmt.Errorf("can't find the section: %w", err)
	}
	return &section, nil
}

func (sectionRepo *sectionRepository) GetSectionsByProfessorAndSemester(ctx context.Context, professorId uint, semesterId uint) ([]models.CourseSection, error) {
	var sections []models.CourseSection
	err := sectionRepo.db.WithContext(ctx).
		Preload("SectionSchedules").
		Where("professor_id = ? AND semester_id = ?", professorId, semesterId).
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch sections for professor: %w", err)
	}
	return sections, nil
}

func (sectionRepo *sectionRepository) GetSectionsByRoomAndSemester(ctx context.Context, roomId uint, semesterId uint) ([]models.CourseSection, error) {
	var sections []models.CourseSection
	err := sectionRepo.db.WithContext(ctx).
		Preload("SectionSchedules").
		Where("room_id = ? AND semester_id = ?", roomId, semesterId).
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch sections for room: %w", err)
	}
	return sections, nil
}

func (sectionRepo *sectionRepository) UpdateSection(ctx context.Context, section *models.CourseSection) error {
	if err := sectionRepo.db.WithContext(ctx).Save(section).Error; err != nil {
		return fmt.Errorf("can't update the section: %w", err)
	}
	return nil
}

func (sectionRepo *sectionRepository) GetSectionsByProfessor(ctx context.Context, professorId uint) ([]models.CourseSection, error) {
	var sections []models.CourseSection
	err := sectionRepo.db.WithContext(ctx).
		Preload("Semester").
		Where("professor_id = ?", professorId).
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("can't fetch sections for the professor: %w", err)
	}
	return sections, nil
}

func (sectionRepo *sectionRepository) SearchSectionsByProfessorAndCourseName(ctx context.Context, professorId uint, name string) ([]models.CourseSection, error) {
	var sections []models.CourseSection
	pattern := "%" + name + "%"

	err := sectionRepo.db.WithContext(ctx).
		Joins("JOIN courses ON courses.id = course_sections.course_id").
		Where("course_sections.professor_id = ? AND courses.course_name ILIKE ?", professorId, pattern).
		Preload("Course").
		Find(&sections).Error
	if err != nil {
		return nil, fmt.Errorf("can't search sections for the professor: %w", err)
	}
	return sections, nil
}

func NewSectionRepository(db *gorm.DB) SectionRepository {
	return &sectionRepository{db: db}
}
