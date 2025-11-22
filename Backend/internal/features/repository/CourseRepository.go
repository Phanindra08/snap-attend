package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

var (
	ErrCourseNotFound = errors.New("course not found")
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course) error
	GetCourseByID(ctx context.Context, id uint) (*models.Course, error)
	GetCourseByName(ctx context.Context, name string) (*models.Course, error)
	UpdateCourse(ctx context.Context, course *models.Course) error
	DeleteCourse(ctx context.Context, id uint) error
	SearchCoursesByName(ctx context.Context, name string) ([]models.Course, error)
}

type courseRepository struct {
	db *gorm.DB
}

func (userRepo *courseRepository) CreateCourse(ctx context.Context, course *models.Course) error {
	if err := userRepo.db.WithContext(ctx).Create(course).Error; err != nil {
		return fmt.Errorf("can't save the course: %w", err)
	}
	return nil
}

func (userRepo *courseRepository) GetCourseByID(ctx context.Context, id uint) (*models.Course, error) {
	var course models.Course
	if err := userRepo.db.WithContext(ctx).First(&course, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, fmt.Errorf("can't find the course: %w", err)
	}
	return &course, nil
}

func (userRepo *courseRepository) GetCourseByName(ctx context.Context, name string) (*models.Course, error) {
	var course models.Course
	if err := userRepo.db.WithContext(ctx).Where("course_name = ?", name).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCourseNotFound
		}
		return nil, fmt.Errorf("can't find the course by name: %w", err)
	}
	return &course, nil
}

func (userRepo *courseRepository) UpdateCourse(ctx context.Context, course *models.Course) error {
	if err := userRepo.db.WithContext(ctx).Save(course).Error; err != nil {
		return fmt.Errorf("can't update the course: %w", err)
	}
	return nil
}

func (userRepo *courseRepository) DeleteCourse(ctx context.Context, id uint) error {
	result := userRepo.db.WithContext(ctx).Delete(&models.Course{}, id)
	if result.Error != nil {
		return fmt.Errorf("can't delete the course: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrCourseNotFound
	}
	return nil
}

func (userRepo *courseRepository) SearchCoursesByName(ctx context.Context, name string) ([]models.Course, error) {
	var courses []models.Course
	pattern := "%" + name + "%"

	if err := userRepo.db.WithContext(ctx).
		Where("course_name ILIKE ?", pattern).
		Find(&courses).Error; err != nil {
		return nil, fmt.Errorf("can't search courses: %w", err)
	}

	return courses, nil
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}
