package service

import (
	"context"
	"errors"
	"time"

	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
)

var (
	ErrProfessorScheduleConflict = errors.New("professor has a schedule conflict with another course section")
)

type SectionService interface {
	AssignProfessorToSection(ctx context.Context, sectionId uint, professorId uint) (*models.CourseSection, error)
}

type sectionService struct {
	sectionRepo repository.SectionRepository
}

func (sectionService *sectionService) AssignProfessorToSection(ctx context.Context, sectionId uint, professorId uint) (*models.CourseSection, error) {
	section, err := sectionService.sectionRepo.GetSectionByID(ctx, sectionId)
	if err != nil {
		return nil, err
	}

	semesterId := section.SemesterId

	professorSections, err := sectionService.sectionRepo.GetSectionsByProfessorAndSemester(ctx, professorId, semesterId)
	if err != nil {
		return nil, err
	}

	for _, other := range professorSections {
		if other.ID == section.ID {
			continue
		}
		for _, newSched := range section.SectionSchedules {
			for _, existingSched := range other.SectionSchedules {
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

	section.ProfessorId = professorId
	if err := sectionService.sectionRepo.UpdateSection(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

func NewSectionService(sectionRepo repository.SectionRepository) SectionService {
	return &sectionService{sectionRepo: sectionRepo}
}
