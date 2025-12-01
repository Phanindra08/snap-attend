package service

import (
	"context"
	"errors"
	"time"

	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
)

var (
	ErrMaxCoursesPerSemester   = errors.New("student has reached maximum number of courses for this semester")
	ErrAlreadyEnrolledInCourse = errors.New("student is already enrolled in this course for the semester")
	ErrStudentScheduleConflict = errors.New("student has a schedule conflict with another course")
	ErrSectionFull             = errors.New("course section is full")
)

type EnrollmentService interface {
	EnrollStudentInSection(ctx context.Context, studentId uint, sectionId uint) (*models.StudentSectionEnrollment, error)
	GetStudentEnrollments(ctx context.Context, studentId uint) ([]models.StudentSectionEnrollment, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	sectionRepo    repository.SectionRepository
}

func (enrollService *enrollmentService) EnrollStudentInSection(ctx context.Context, studentId uint, sectionId uint) (*models.StudentSectionEnrollment, error) {
	section, err := enrollService.sectionRepo.GetSectionByID(ctx, sectionId)
	if err != nil {
		return nil, err
	}

	semesterId := section.SemesterId

	count, err := enrollService.enrollmentRepo.CountEnrollmentsByStudentAndSemester(ctx, studentId, semesterId)
	if err != nil {
		return nil, err
	}
	if count >= 4 {
		return nil, ErrMaxCoursesPerSemester
	}

	enrollments, err := enrollService.enrollmentRepo.GetEnrollmentsByStudentAndSemester(ctx, studentId, semesterId)
	if err != nil {
		return nil, err
	}

	for _, e := range enrollments {
		existingSection := e.CourseSection

		if existingSection.CourseId == section.CourseId {
			return nil, ErrAlreadyEnrolledInCourse
		}

		for _, newSched := range section.SectionSchedules {
			for _, existingSched := range existingSection.SectionSchedules {
				if newSched.DaysOfTheClass != existingSched.DaysOfTheClass {
					continue
				}

				newStartStr := newSched.StartTime.String()
				newEndStr := newSched.EndTime.String()
				existingStartStr := existingSched.StartTime.String()
				existingEndStr := existingSched.EndTime.String()

				newStart, err := time.Parse("15:04:05", newStartStr)
				if err != nil {
					return nil, err
				}
				newEnd, err := time.Parse("15:04:05", newEndStr)
				if err != nil {
					return nil, err
				}
				existingStart, err := time.Parse("15:04:05", existingStartStr)
				if err != nil {
					return nil, err
				}
				existingEnd, err := time.Parse("15:04:05", existingEndStr)
				if err != nil {
					return nil, err
				}

				if newStart.Before(existingEnd) && newEnd.After(existingStart) {
					return nil, ErrStudentScheduleConflict
				}
			}
		}
	}

	if section.EnrolledSeats >= section.TotalSeats {
		return nil, ErrSectionFull
	}

	section.EnrolledSeats++
	if err := enrollService.sectionRepo.UpdateSection(ctx, section); err != nil {
		return nil, err
	}

	enrollment := &models.StudentSectionEnrollment{
		StudentId:  studentId,
		SectionId:  sectionId,
		EnrolledAt: time.Now().UTC(),
	}

	if err := enrollService.enrollmentRepo.CreateEnrollment(ctx, enrollment); err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (enrollService *enrollmentService) GetStudentEnrollments(ctx context.Context, studentId uint) ([]models.StudentSectionEnrollment, error) {
	enrollments, err := enrollService.enrollmentRepo.GetEnrollmentsByStudent(ctx, studentId)
	if err != nil {
		return nil, err
	}

	var activeEnrollments []models.StudentSectionEnrollment
	now := time.Now().UTC()

	for _, enrollment := range enrollments {
		semester := enrollment.CourseSection.Semester
		
		// Converting datatypes.Date to time.Time
		startTime := time.Time(semester.StartDate)
		endTime := time.Time(semester.EndDate)

		startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
		endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, time.UTC)

		if (now.Equal(startTime) || now.After(startTime)) && (now.Equal(endTime) || now.Before(endTime)) {
			activeEnrollments = append(activeEnrollments, enrollment)
		}
	}

	return activeEnrollments, nil
}

func NewEnrollmentService(enrollmentRepo repository.EnrollmentRepository, sectionRepo repository.SectionRepository) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		sectionRepo:    sectionRepo,
	}
}
