package service

import (
	"context"
	"errors"
	"strings"

	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

var (
	ErrCourseAlreadyExists = errors.New("course already exists")
)

type CourseService interface {
	CreateCourse(ctx context.Context, courseName string) (*models.Course, error)
	GetCourseByID(ctx context.Context, id uint) (*models.Course, error)
	UpdateCourse(ctx context.Context, id uint, courseName string) (*models.Course, error)
	DeleteCourse(ctx context.Context, id uint) error
	SearchCoursesByName(ctx context.Context, name string) ([]models.Course, error)
}

type courseService struct {
	courseRepo repository.CourseRepository
}

func (courseService *courseService) CreateCourse(ctx context.Context, courseName string) (*models.Course, error) {
	name := strings.TrimSpace(courseName)
	if utils.IsEmpty(name) {
		return nil, errors.New("course name cannot be empty")
	}

	// Checking if the course already exists
	_, err := courseService.courseRepo.GetCourseByName(ctx, name)
	if err == nil {
		return nil, ErrCourseAlreadyExists
	}
	if !errors.Is(err, repository.ErrCourseNotFound) {
		return nil, err
	}

	course := &models.Course{
		CourseName: name,
	}

	if err := courseService.courseRepo.CreateCourse(ctx, course); err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") {
			return nil, ErrCourseAlreadyExists
		}
		return nil, err
	}

	return course, nil
}

func (courseService *courseService) GetCourseByID(ctx context.Context, id uint) (*models.Course, error) {
	course, err := courseService.courseRepo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (courseService *courseService) UpdateCourse(ctx context.Context, id uint, courseName string) (*models.Course, error) {
	course, err := courseService.courseRepo.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(courseName)
	if utils.IsEmpty(name) {
		return nil, errors.New("course name cannot be empty")
	}

	course.CourseName = name
	if err := courseService.courseRepo.UpdateCourse(ctx, course); err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") {
			return nil, ErrCourseAlreadyExists
		}
		return nil, err
	}

	return course, nil
}

func (courseService *courseService) DeleteCourse(ctx context.Context, id uint) error {
	if err := courseService.courseRepo.DeleteCourse(ctx, id); err != nil {
		return err
	}
	return nil
}

func (courseService *courseService) SearchCoursesByName(ctx context.Context, name string) ([]models.Course, error) {
	trimmedCourseName := strings.TrimSpace(name)
	if utils.IsEmpty(trimmedCourseName) {
		return nil, errors.New("course name cannot be empty")
	}
	return courseService.courseRepo.SearchCoursesByName(ctx, trimmedCourseName)
}

func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{courseRepo: courseRepo}
}
