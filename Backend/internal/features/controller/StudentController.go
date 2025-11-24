package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanindra08/snap-attend/internal/features/InputRequest"
	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

type StudentController struct {
	attendanceService service.AttendanceService
}

func (sc *StudentController) SubmitAttendance(ctx *gin.Context) {
	studentID := ctx.GetUint("userID")

	var attendanceRequest InputRequest.StudentSubmitAttendanceRequest
	if err := ctx.ShouldBindJSON(&attendanceRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trimming whitespaces at the beginning and end of input fields
	attendanceRequest.Question = utils.TrimAndConvertToLowerCase(attendanceRequest.Question)
	attendanceRequest.Answer = strings.TrimSpace(attendanceRequest.Answer)

	input := dto.StudentSubmitAttendanceDTO{
		QrId:      attendanceRequest.QrId,
		Latitude:  attendanceRequest.Latitude,
		Longitude: attendanceRequest.Longitude,
		Question:  attendanceRequest.Question,
		Answer:    attendanceRequest.Answer,
	}

	attendance, err := sc.attendanceService.SubmitAttendance(
		ctx.Request.Context(),
		studentID,
		input,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAttendanceLocation):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "You are not at the class location"})
			return
		case errors.Is(err, service.ErrStudentNotEnrolledInSection):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student is not enrolled in this course"})
			return
		case errors.Is(err, service.ErrAttendanceQrInactive):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Attendance QR is inactive"})
			return
		case errors.Is(err, service.ErrAttendanceQrExpired):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Attendance QR is expired"})
			return
		case errors.Is(err, service.ErrAttendanceQrNotYetValid):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Attendance QR is not yet valid"})
			return
		case errors.Is(err, service.ErrAttendanceAlreadySubmitted):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Attendance has already been submitted for this QR"})
			return
		case errors.Is(err, service.ErrClassNotInSession):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Class is not in session at this time"})
			return
		case errors.Is(err, service.ErrSectionNotActive):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course section or semester is not active"})
			return
		case errors.Is(err, service.ErrInvalidLocationData):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing location data"})
			return
		case errors.Is(err, repository.ErrAttendanceQrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Attendance QR not found"})
			return
		default:
			log.Printf("Student submit attendance error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit attendance"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Attendance submitted successfully",
		"attendance": gin.H{
			"id":         attendance.ID,
			"studentId":  attendance.StudentId,
			"qrId":       attendance.QrId,
			"attendedAt": attendance.AttendedAt,
			"attended":   attendance.Attended,
		},
	})
}

func (sc *StudentController) GetAttendanceSummary(ctx *gin.Context) {
	studentID := ctx.GetUint("userID")

	courseIdInStringFormat := ctx.Query("courseId")
	if utils.IsEmpty(courseIdInStringFormat) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "courseId query parameter is required"})
		return
	}

	courseIdUint64, err := strconv.ParseUint(courseIdInStringFormat, 10, 64)
	if err != nil || courseIdUint64 == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid courseId"})
		return
	}

	summaries, err := sc.attendanceService.GetAttendanceSummary(ctx.Request.Context(), studentID, uint(courseIdUint64))
	if err != nil {
		log.Printf("Student attendance summary error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance summary"})
		return
	}

	var result []gin.H
	for _, summary := range summaries {
		result = append(result, gin.H{
			"courseId":        summary.CourseID,
			"courseName":      summary.CourseName,
			"sectionId":       summary.SectionID,
			"sectionNumber":   summary.SectionNumber,
			"attendedClasses": summary.AttendedCount,
			"totalClasses":    summary.TotalSessions,
			"attendancePct":   summary.AttendancePct,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"summary": result,
	})
}

func (sc *StudentController) GetAttendanceHistory(ctx *gin.Context) {
	studentID := ctx.GetUint("userID")

	startDateInStringFormat := ctx.Query("startDate")
	endDateInStringFormat := ctx.Query("endDate")

	var fromDate *time.Time
	var toDate *time.Time

	if utils.IsEmpty(startDateInStringFormat) {
		startDate, err := time.Parse("2006-01-02", startDateInStringFormat)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD"})
			return
		}
		fromDate = &startDate
	}

	if utils.IsEmpty(endDateInStringFormat) {
		endDate, err := time.Parse("2006-01-02", endDateInStringFormat)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
			return
		}
		toDate = &endDate
	}

	records, err := sc.attendanceService.GetAttendanceHistory(ctx.Request.Context(), studentID, fromDate, toDate)
	if err != nil {
		log.Printf("Student attendance history error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance history"})
		return
	}

	var history []gin.H
	for _, record := range records {
		section := record.AttendanceQr.CourseSection
		course := section.Course

		history = append(history, gin.H{
			"id":            record.ID,
			"qrId":          record.QrId,
			"courseId":      course.ID,
			"courseName":    course.CourseName,
			"sectionId":     section.ID,
			"sectionNumber": section.SectionNumber,
			"attendedAt":    record.AttendedAt,
			"attended":      record.Attended,
			"question":      record.StudentQuestion,
			"answer":        record.StudentAnswer,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

func NewStudentController(attendanceService service.AttendanceService) *StudentController {
	return &StudentController{
		attendanceService: attendanceService,
	}
}
