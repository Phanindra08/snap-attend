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
	CreateSection(ctx context.Context, courseId uint, sectionNumber string, semesterId uint, roomId uint, latitude float64, longitude float64) (*models.CourseSection, error)
	GetAllSections(ctx context.Context) ([]models.CourseSection, error)
	UpdateSection(ctx context.Context, sectionId uint, sectionNumber *string, latitude *float64, longitude *float64, roomId *uint) (*models.CourseSection, error)
}

type sectionService struct {
	sectionRepo repository.SectionRepository
	roomRepo    repository.RoomRepository
}

func (sectionService *sectionService) CreateSection(ctx context.Context, courseId uint, sectionNumber string, semesterId uint, roomId uint, latitude float64, longitude float64) (*models.CourseSection, error) {
	room, err := sectionService.roomRepo.GetRoomByID(ctx, roomId)
	if err != nil {
		return nil, err
	}

	section := &models.CourseSection{
		CourseId:      courseId,
		SectionNumber: sectionNumber,
		SemesterId:    semesterId,
		RoomId:        roomId,
		TotalSeats:    uint(room.Capacity),
		EnrolledSeats: 0,
		Latitude:      latitude,
		Longitude:     longitude,
	}

	if err := sectionService.sectionRepo.CreateSection(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
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

func (sectionService *sectionService) GetAllSections(ctx context.Context) ([]models.CourseSection, error) {
	return sectionService.sectionRepo.GetAllSections(ctx)
}

func (sectionService *sectionService) UpdateSection(ctx context.Context, sectionId uint, sectionNumber *string, latitude *float64, longitude *float64, roomId *uint) (*models.CourseSection, error) {
	section, err := sectionService.sectionRepo.GetSectionByID(ctx, sectionId)
	if err != nil {
		return nil, err
	}

	if sectionNumber != nil {
		section.SectionNumber = *sectionNumber
	}
	if latitude != nil {
		section.Latitude = *latitude
	}
	if longitude != nil {
		section.Longitude = *longitude
	}
	if roomId != nil {
		section.RoomId = *roomId
	}

	if err := sectionService.sectionRepo.UpdateSection(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

func NewSectionService(sectionRepo repository.SectionRepository, roomRepo repository.RoomRepository) SectionService {
	return &sectionService{
		sectionRepo: sectionRepo,
		roomRepo:    roomRepo,
	}
}
