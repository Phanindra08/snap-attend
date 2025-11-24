package service

import (
	"context"
	"errors"
	"time"

	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

var (
	ErrInvalidAttendanceLocation   = errors.New("student is not at the class location")
	ErrStudentNotEnrolledInSection = errors.New("student is not enrolled in this course section")
	ErrAttendanceQrInactive        = errors.New("attendance QR is inactive")
	ErrAttendanceQrExpired         = errors.New("attendance QR is expired")
	ErrAttendanceQrNotYetValid     = errors.New("attendance QR is not yet valid")
	ErrAttendanceAlreadySubmitted  = errors.New("attendance already submitted for this QR")
	ErrClassNotInSession           = errors.New("class is not in the session at this time")
	ErrSectionNotActive            = errors.New("course section / semester is not active")
	ErrInvalidLocationData         = errors.New("invalid or missing location data")
)

type AttendanceSummary struct {
	CourseID      uint
	CourseName    string
	SectionID     uint
	SectionNumber string
	AttendedCount int64
	TotalSessions int64
	AttendancePct float64
}

type AttendanceService interface {
	SubmitAttendance(ctx context.Context, studentID uint, attendanceDTO dto.StudentSubmitAttendanceDTO) (*models.StudentAttendance, error)
	GetAttendanceSummary(ctx context.Context, studentId uint, courseId uint) ([]AttendanceSummary, error)
	GetAttendanceHistory(ctx context.Context, studentId uint, from *time.Time, to *time.Time) ([]models.StudentAttendance, error)
}

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
	enrollRepo     repository.EnrollmentRepository
}

func (attendanceService *attendanceService) SubmitAttendance(ctx context.Context, studentID uint, attendanceDTO dto.StudentSubmitAttendanceDTO) (*models.StudentAttendance, error) {
	qr, err := attendanceService.attendanceRepo.GetAttendanceQrWithSectionAndRoom(ctx, attendanceDTO.QrId)
	if err != nil {
		return nil, err
	}

	section := qr.CourseSection
	room := section.Room

	now := time.Now().UTC()

	// Check that QR is active and valid (based on GeneratedAt and ExpiresAt)
	if !qr.IsActive {
		return nil, ErrAttendanceQrInactive
	}
	if now.Before(qr.GeneratedAt) {
		return nil, ErrAttendanceQrNotYetValid
	}
	if now.After(qr.ExpiresAt) {
		return nil, ErrAttendanceQrExpired
	}

	// 2) Check that the section / semester is active (based on Semester dates)
	semester := section.Semester
	startDate := semester.StartDate
	endDate := semester.EndDate

	// Converting to time.Time
	start := time.Time(startDate)
	end := time.Time(endDate)

	// Normalize to day boundaries (inclusive)
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)

	if now.Before(start) || now.After(end) {
		return nil, ErrSectionNotActive
	}

	// Check class is aligned with SectionSchedules + day of week + time
	if !utils.IsClassInSessionNow(now, section.SectionSchedules) {
		return nil, ErrClassNotInSession
	}

	// Checking if student is enrolled in the section
	enrollment, err := attendanceService.enrollRepo.GetEnrollmentByStudentAndSection(ctx, studentID, section.ID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, ErrStudentNotEnrolledInSection
	}

	// Check if attendance was already submitted for this QR by the student
	existingAttendance, err := attendanceService.attendanceRepo.GetStudentAttendanceByStudentAndQr(ctx, studentID, qr.ID)
	if err != nil {
		return nil, err
	}
	if existingAttendance != nil {
		return nil, ErrAttendanceAlreadySubmitted
	}

	// Validating the student location
	if !utils.IsValidLatAndLon(attendanceDTO.Latitude, attendanceDTO.Longitude) {
		return nil, ErrInvalidLocationData
	}

	// Validate the student location against the room coordinates
	distance := utils.CalculateDistanceInMeters(
		attendanceDTO.Latitude,
		attendanceDTO.Longitude,
		room.Latitude,
		room.Longitude,
	)

	if distance > utils.MAX_ATTENDANCE_DISTANCE_METERS {
		return nil, ErrInvalidAttendanceLocation
	}

	attendance := &models.StudentAttendance{
		StudentId:       studentID,
		QrId:            qr.ID,
		AttendedAt:      now,
		Attended:        true,
		StudentQuestion: attendanceDTO.Question,
		StudentAnswer:   attendanceDTO.Answer,
	}

	if err := attendanceService.attendanceRepo.CreateStudentAttendance(ctx, attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

func (attendanceService *attendanceService) GetAttendanceSummary(ctx context.Context, studentId uint, courseId uint) ([]AttendanceSummary, error) {
	attendanceRecords, err := attendanceService.attendanceRepo.GetStudentAttendancesByStudent(ctx, studentId)
	if err != nil {
		return nil, err
	}

	type agg struct {
		courseID      uint
		courseName    string
		sectionID     uint
		sectionNumber string
		attended      int64
		totalSessions int64
	}

	aggregates := map[uint]*agg{}

	for _, rec := range attendanceRecords {
		section := rec.AttendanceQr.CourseSection
		course := section.Course

		if course.ID != courseId {
			continue
		}

		sectionID := section.ID

		if _, exists := aggregates[sectionID]; !exists {
			totalSessions, err := attendanceService.attendanceRepo.CountAttendanceQrBySectionUntilNow(ctx, sectionID)
			if err != nil {
				return nil, err
			}

			aggregates[sectionID] = &agg{
				courseID:      course.ID,
				courseName:    course.CourseName,
				sectionID:     sectionID,
				sectionNumber: section.SectionNumber,
				attended:      0,
				totalSessions: totalSessions,
			}
		}

		if rec.Attended {
			aggregates[sectionID].attended++
		}
	}

	var attendanceSummary []AttendanceSummary
	for _, a := range aggregates {
		var percentage float64
		if a.totalSessions > 0 {
			percentage = (float64(a.attended) / float64(a.totalSessions)) * 100.0
		}

		attendanceSummary = append(attendanceSummary, AttendanceSummary{
			CourseID:      a.courseID,
			CourseName:    a.courseName,
			SectionID:     a.sectionID,
			SectionNumber: a.sectionNumber,
			AttendedCount: a.attended,
			TotalSessions: a.totalSessions,
			AttendancePct: percentage,
		})
	}

	return attendanceSummary, nil
}

func (attendanceService *attendanceService) GetAttendanceHistory(ctx context.Context, studentId uint, from *time.Time, to *time.Time) ([]models.StudentAttendance, error) {
	return attendanceService.attendanceRepo.GetStudentAttendanceHistory(ctx, studentId, from, to)
}

func NewAttendanceService(attendanceRepo repository.AttendanceRepository, enrollRepo repository.EnrollmentRepository) AttendanceService {
	return &attendanceService{
		attendanceRepo: attendanceRepo,
		enrollRepo:     enrollRepo,
	}
}
