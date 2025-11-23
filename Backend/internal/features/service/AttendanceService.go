package service

import (
	"context"
	"errors"
	"time"

	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

var (
	ErrInvalidAttendanceLocation   = errors.New("student is not at the class location")
	ErrStudentNotEnrolledInSection = errors.New("student is not enrolled in this course section")
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
	SubmitAttendance(ctx context.Context, studentId uint, qrId uint, latitude float64, longitude float64, question string, answer string) (*models.StudentAttendance, error)
	GetAttendanceSummary(ctx context.Context, studentId uint, courseId uint) ([]AttendanceSummary, error)
	GetAttendanceHistory(ctx context.Context, studentId uint, from *time.Time, to *time.Time) ([]models.StudentAttendance, error)
}

type attendanceService struct {
	attendanceRepo repository.AttendanceRepository
	enrollRepo     repository.EnrollmentRepository
}

func (attendanceService *attendanceService) SubmitAttendance(ctx context.Context, studentId uint, qrId uint, latitude float64, longitude float64, question string, answer string) (*models.StudentAttendance, error) {
	qr, err := attendanceService.attendanceRepo.GetAttendanceQrWithSectionAndRoom(ctx, qrId)
	if err != nil {
		return nil, err
	}

	section := qr.CourseSection
	room := section.Room

	// Check if student is enrolled in the section
	enrollment, err := attendanceService.enrollRepo.GetEnrollmentByStudentAndSection(ctx, studentId, section.ID)
	if err != nil {
		return nil, err
	}
	if enrollment == nil {
		return nil, ErrStudentNotEnrolledInSection
	}

	// Validate the student location against the room coordinates
	distance := utils.CalculateDistanceInMeters(
		latitude,
		longitude,
		room.Latitude,
		room.Longitude,
	)

	if distance > utils.MAX_ATTENDANCE_DISTANCE_METERS {
		return nil, ErrInvalidAttendanceLocation
	}

	attendance := &models.StudentAttendance{
		StudentId:       studentId,
		QrId:            qr.ID,
		AttendedAt:      time.Now().UTC(),
		Attended:        true,
		StudentQuestion: question,
		StudentAnswer:   answer,
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
