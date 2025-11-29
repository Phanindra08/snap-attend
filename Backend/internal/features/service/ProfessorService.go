package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

var (
	ErrProfessorNotOwner    = errors.New("this section is not offered by you. Please select the course which you are teaching")
	ErrProfessorSearchEmpty = errors.New("at least one search parameter must be provided")
	ErrAttendanceNotFound   = errors.New("attendance record not found")
)

type ProfessorService interface {
	GenerateAttendanceQr(ctx context.Context, professorID uint, sectionID uint) (*models.AttendanceQr, error)
	GetDailyAttendanceReport(ctx context.Context, professorID uint, sectionID uint, day time.Time) ([]dto.DailyAttendanceRow, error)
	GetSectionAttendance(ctx context.Context, professorID uint, sectionID uint, fromDate *time.Time, toDate *time.Time, studentId *uint) ([]repository.StudentAttendanceWithUser, error)
	UpdateAttendanceStatus(ctx context.Context, professorID uint, attendanceID uint, attended bool) (*models.StudentAttendance, error)
	Search(ctx context.Context, professorID uint, courseName string, studentName string) (*dto.ProfessorSearchResult, error)
	GetSectionAttendanceOverview(ctx context.Context, professorID uint, sectionID uint) (*dto.SectionAttendanceOverview, error)
}

type professorService struct {
	sectionRepo    repository.SectionRepository
	enrollmentRepo repository.EnrollmentRepository
	attendanceRepo repository.AttendanceRepository
}

func (ps *professorService) GenerateAttendanceQr(ctx context.Context, professorID uint, sectionID uint) (*models.AttendanceQr, error) {
	section, err := ps.sectionRepo.GetSectionByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section.ProfessorId != professorID {
		return nil, ErrProfessorNotOwner
	}

	now := time.Now().UTC()
	if !utils.IsClassInSessionNow(now, section.SectionSchedules) {
		return nil, errors.New("class is not currently in the scheduled time")
	}

	// Deactivating old active QRs for this section
	if err := ps.attendanceRepo.DeactivateActiveQrsForSection(ctx, sectionID); err != nil {
		// Not a major error. This is just a cleanup step. Log and continue.
		log.Printf("WARNING: failed to deactivate previous QRs for section %d: %v", sectionID, err)
	}

	generatedAt := now
	expiresAt := generatedAt.Add(time.Minute * utils.QR_TTL_MINUTES)

	newQr := &repository.AttendanceQrModelForCreate{
		SectionId:   sectionID,
		GeneratedAt: generatedAt,
		ExpiresAt:   expiresAt,
	}

	createdQr, err := ps.attendanceRepo.CreateAttendanceQr(ctx, newQr)
	if err != nil {
		return nil, err
	}

	return createdQr, nil
}

func (ps *professorService) GetDailyAttendanceReport(ctx context.Context, professorID uint, sectionID uint, day time.Time) ([]dto.DailyAttendanceRow, error) {
	section, err := ps.sectionRepo.GetSectionByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section.ProfessorId != professorID {
		return nil, ErrProfessorNotOwner
	}

	fromDate := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	toDate := fromDate.Add(24*time.Hour - time.Nanosecond)

	studentAttendanceBySection, err := ps.attendanceRepo.GetStudentAttendanceBySection(ctx, sectionID, &fromDate, &toDate, nil)
	if err != nil {
		return nil, err
	}

	enrollments, err := ps.enrollmentRepo.GetEnrollmentsBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	attendanceByStudent := make(map[uint][]repository.StudentAttendanceWithUser)
	for _, record := range studentAttendanceBySection {
		attendanceByStudent[record.StudentId] = append(attendanceByStudent[record.StudentId], record)
	}

	var dailyAttendanceRows []dto.DailyAttendanceRow
	for _, enrollment := range enrollments {
		student := enrollment.Student
		attendanceRecords := attendanceByStudent[student.ID]

		present := false
		var attendedAt *time.Time

		for _, attendanceRecord := range attendanceRecords {
			if attendanceRecord.Attended {
				present = true
				attendedTime := attendanceRecord.AttendedAt
				attendedAt = &attendedTime
				break
			}
		}

		dailyAttendanceRows = append(dailyAttendanceRows, dto.DailyAttendanceRow{
			Student:    student,
			Present:    present,
			AttendedAt: attendedAt,
		})
	}

	return dailyAttendanceRows, nil
}

func (ps *professorService) GetSectionAttendance(ctx context.Context, professorID uint, sectionID uint, fromDate *time.Time, toDate *time.Time, studentId *uint) ([]repository.StudentAttendanceWithUser, error) {
	section, err := ps.sectionRepo.GetSectionByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section.ProfessorId != professorID {
		return nil, ErrProfessorNotOwner
	}

	return ps.attendanceRepo.GetStudentAttendanceBySection(ctx, sectionID, fromDate, toDate, studentId)
}

func (ps *professorService) UpdateAttendanceStatus(ctx context.Context, professorID uint, attendanceID uint, attended bool) (*models.StudentAttendance, error) {
	record, err := ps.attendanceRepo.GetAttendanceByIDWithSection(ctx, attendanceID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrAttendanceNotFound
	}

	section := record.AttendanceQr.CourseSection
	if section.ProfessorId != professorID {
		return nil, ErrProfessorNotOwner
	}

	record.Attended = attended

	if err := ps.attendanceRepo.UpdateStudentAttendance(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}

func (ps *professorService) Search(ctx context.Context, professorID uint, courseName string, studentName string) (*dto.ProfessorSearchResult, error) {
	if utils.IsEmpty(courseName) && utils.IsEmpty(studentName) {
		return nil, ErrProfessorSearchEmpty
	}

	result := &dto.ProfessorSearchResult{}

	// Searching by Course name
	if !utils.IsEmpty(courseName) {
		sections, err := ps.sectionRepo.SearchSectionsByProfessorAndCourseName(ctx, professorID, courseName)
		if err != nil {
			return nil, err
		}

		for _, section := range sections {
			result.Courses = append(result.Courses, dto.CourseSearchResult{
				CourseID:      section.Course.ID,
				CourseName:    section.Course.CourseName,
				SectionID:     section.ID,
				SectionNumber: section.SectionNumber,
			})
		}
	}

	// Searching by Student name (current semester only)
	if !utils.IsEmpty(studentName) {
		// Convert the search term into lowercase and trim spaces
		searchName := utils.TrimAndConvertToLowerCase(studentName)

		// Getting all the sections for professor, but with Semester preloaded
		sections, err := ps.sectionRepo.GetSectionsByProfessor(ctx, professorID)
		if err != nil {
			return nil, err
		}

		now := time.Now().UTC()
		seen := make(map[uint]bool)

		for _, section := range sections {
			sem := section.Semester

			// Converting datatypes.Date to time.Time
			startTime := time.Time(sem.StartDate)
			endTime := time.Time(sem.EndDate)

			startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, time.UTC)
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 0, time.UTC)

			// Considering the sections where now is within the semester start or end date
			if now.Before(startTime) || now.After(endTime) {
				continue
			}

			// Getting enrollments for the current semester
			enrollments, err := ps.enrollmentRepo.GetEnrollmentsBySection(ctx, section.ID)
			if err != nil {
				return nil, err
			}

			for _, enrollment := range enrollments {
				student := enrollment.Student
				if seen[student.ID] {
					continue
				}
				fullName := utils.TrimAndConvertToLowerCase(student.FirstName + " " + student.LastName)
				if strings.Contains(fullName, searchName) {
					seen[student.ID] = true
					result.Students = append(result.Students, dto.StudentSearchResult{
						StudentID: student.ID,
						FirstName: student.FirstName,
						LastName:  student.LastName,
						Email:     student.Email,
					})
				}
			}
		}
	}

	return result, nil
}

func (ps *professorService) GetSectionAttendanceOverview(ctx context.Context, professorID uint, sectionID uint) (*dto.SectionAttendanceOverview, error) {
	section, err := ps.sectionRepo.GetSectionByID(ctx, sectionID)
	if err != nil {
		return nil, err
	}
	if section.ProfessorId != professorID {
		return nil, ErrProfessorNotOwner
	}

	totalSessions, err := ps.attendanceRepo.CountAttendanceQrBySectionUntilNow(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	enrollments, err := ps.enrollmentRepo.GetEnrollmentsBySection(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	overview := &dto.SectionAttendanceOverview{
		SectionID:     sectionID,
		TotalSessions: totalSessions,
	}

	for _, enrollment := range enrollments {
		student := enrollment.Student

		attendedCount, err := ps.attendanceRepo.CountStudentAttendanceBySection(ctx, student.ID, sectionID)
		if err != nil {
			return nil, err
		}

		var pct float64
		if totalSessions > 0 {
			pct = (float64(attendedCount) / float64(totalSessions)) * 100.0
		}

		overview.Students = append(overview.Students, dto.SectionAttendanceOverviewRow{
			StudentID:       student.ID,
			FirstName:       student.FirstName,
			LastName:        student.LastName,
			Email:           student.Email,
			AttendedClasses: attendedCount,
			TotalClasses:    totalSessions,
			AttendancePct:   pct,
		})
	}

	return overview, nil
}

func NewProfessorService(
	sectionRepo repository.SectionRepository,
	enrollmentRepo repository.EnrollmentRepository,
	attendanceRepo repository.AttendanceRepository,
) ProfessorService {
	return &professorService{
		sectionRepo:    sectionRepo,
		enrollmentRepo: enrollmentRepo,
		attendanceRepo: attendanceRepo,
	}
}
