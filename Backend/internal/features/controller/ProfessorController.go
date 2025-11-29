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
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

type ProfessorController struct {
	professorService service.ProfessorService
}

func (pc *ProfessorController) GenerateAttendanceQr(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	sectionID, err := parseUintParam(ctx, "sectionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sectionId"})
		return
	}

	qr, err := pc.professorService.GenerateAttendanceQr(ctx.Request.Context(), professorID, sectionID)
	if err != nil {
		if errors.Is(err, service.ErrProfessorNotOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
			return
		}
		log.Printf("error while generating QR by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate attendance QR"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Attendance QR generated successfully",
		"qr": gin.H{
			"id":          qr.ID,
			"sectionId":   qr.SectionId,
			"generatedAt": qr.GeneratedAt,
			"expiresAt":   qr.ExpiresAt,
			"qrHint":      qr.QrLink,
		},
	})
}

func (pc *ProfessorController) GetDailyAttendanceReport(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	sectionID, err := parseUintParam(ctx, "sectionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sectionId"})
		return
	}

	dateInStringFormat := strings.TrimSpace(ctx.Query("date"))
	if utils.IsEmpty(dateInStringFormat) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date query parameter is required (YYYY-MM-DD)"})
		return
	}

	day, err := time.Parse("2006-01-02", dateInStringFormat)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	attendanceReport, err := pc.professorService.GetDailyAttendanceReport(ctx.Request.Context(), professorID, sectionID, day)
	if err != nil {
		if errors.Is(err, service.ErrProfessorNotOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
			return
		}
		log.Printf("Error while fetching daily attendance report by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	var result []gin.H
	for _, row := range attendanceReport {
		result = append(result, gin.H{
			"studentId":  row.Student.ID,
			"firstName":  row.Student.FirstName,
			"lastName":   row.Student.LastName,
			"email":      row.Student.Email,
			"present":    row.Present,
			"attendedAt": row.AttendedAt,
			"question":   nil,
			"answer":     nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"date":     dateInStringFormat,
		"section":  sectionID,
		"students": result,
	})
}

func (pc *ProfessorController) GetSectionAttendance(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	sectionID, err := parseUintParam(ctx, "sectionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sectionId"})
		return
	}

	var (
		fromDate  *time.Time
		toDate    *time.Time
		studentId *uint
	)

	startDateInStringFormat := strings.TrimSpace(ctx.Query("startDate"))
	if !utils.IsEmpty(startDateInStringFormat) {
		startDate, err := time.Parse("2006-01-02", startDateInStringFormat)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid startDate format. Use YYYY-MM-DD"})
			return
		}
		formattedStartDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		fromDate = &formattedStartDate
	}

	endDateInStringFormat := strings.TrimSpace(ctx.Query("endDate"))
	if !utils.IsEmpty(endDateInStringFormat) {
		endDate, err := time.Parse("2006-01-02", endDateInStringFormat)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
			return
		}
		formattedEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, time.UTC)
		toDate = &formattedEndDate
	}

	studentIdInStringFormat := strings.TrimSpace(ctx.Query("studentId"))
	if !utils.IsEmpty(studentIdInStringFormat) {
		id64, err := strconv.ParseUint(studentIdInStringFormat, 10, 64)
		if err != nil || id64 == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid studentId"})
			return
		}
		formattedStudentId := uint(id64)
		studentId = &formattedStudentId
	}

	records, err := pc.professorService.GetSectionAttendance(ctx.Request.Context(), professorID, sectionID, fromDate, toDate, studentId)
	if err != nil {
		if errors.Is(err, service.ErrProfessorNotOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
			return
		}
		log.Printf("Error while querying the sections attendance by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	var result []gin.H
	for _, record := range records {
		student := record.User
		section := record.AttendanceQr.CourseSection
		course := section.Course

		result = append(result, gin.H{
			"attendanceId":  record.ID,
			"studentId":     student.ID,
			"firstName":     student.FirstName,
			"lastName":      student.LastName,
			"email":         student.Email,
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

	ctx.JSON(http.StatusOK, gin.H{"attendance": result})
}

func (pc *ProfessorController) UpdateAttendanceStatus(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	attendanceID, err := parseUintParam(ctx, "attendanceId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance record"})
		return
	}

	var attendanceStatusRequest InputRequest.UpdateAttendanceStatusRequest
	if err := ctx.ShouldBindJSON(&attendanceStatusRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if attendanceStatusRequest.Attended == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while submitting the form"})
		return
	}

	record, err := pc.professorService.UpdateAttendanceStatus(
		ctx.Request.Context(),
		professorID,
		attendanceID,
		*attendanceStatusRequest.Attended,
	)
	if err != nil {
		if errors.Is(err, service.ErrProfessorNotOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
			return
		}
		if errors.Is(err, service.ErrAttendanceNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Attendance record is not found"})
			return
		}
		log.Printf("Error while updating the attendance by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Attendance updated successfully",
		"attendanceId": record.ID,
		"attended":     record.Attended,
	})
}

func (pc *ProfessorController) Search(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	courseName := strings.TrimSpace(ctx.Query("courseName"))
	studentName := strings.TrimSpace(ctx.Query("studentName"))

	professorSearchResult, err := pc.professorService.Search(ctx.Request.Context(), professorID, courseName, studentName)
	if err != nil {
		if errors.Is(err, service.ErrProfessorSearchEmpty) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter must be provided"})
			return
		}
		log.Printf("Error while searching by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	response := gin.H{}

	if len(professorSearchResult.Courses) > 0 {
		var courseResults []gin.H
		for _, searchResult := range professorSearchResult.Courses {
			courseResults = append(courseResults, gin.H{
				"courseId":      searchResult.CourseID,
				"courseName":    searchResult.CourseName,
				"sectionId":     searchResult.SectionID,
				"sectionNumber": searchResult.SectionNumber,
			})
		}
		response["courses"] = courseResults
	}

	if len(professorSearchResult.Students) > 0 {
		var students []gin.H
		for _, searchResult := range professorSearchResult.Students {
			students = append(students, gin.H{
				"studentId": searchResult.StudentID,
				"firstName": searchResult.FirstName,
				"lastName":  searchResult.LastName,
				"email":     searchResult.Email,
			})
		}
		response["students"] = students
	}

	ctx.JSON(http.StatusOK, response)
}

func (pc *ProfessorController) GetSectionAttendanceOverview(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	sectionID, err := parseUintParam(ctx, "sectionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sectionId"})
		return
	}

	overview, err := pc.professorService.GetSectionAttendanceOverview(ctx.Request.Context(), professorID, sectionID)
	if err != nil {
		if errors.Is(err, service.ErrProfessorNotOwner) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
			return
		}
		log.Printf("Error while retrieving section attendance overview by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute attendance"})
		return
	}

	var result []gin.H
	for _, row := range overview.Students {
		result = append(result, gin.H{
			"studentId":       row.StudentID,
			"firstName":       row.FirstName,
			"lastName":        row.LastName,
			"email":           row.Email,
			"attendedClasses": row.AttendedClasses,
			"totalClasses":    row.TotalClasses,
			"attendancePct":   row.AttendancePct,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sectionId":          overview.SectionID,
		"totalSessions":      overview.TotalSessions,
		"studentsAttendance": result,
	})
}

func parseUintParam(ctx *gin.Context, name string) (uint, error) {
	val := ctx.Param(name)
	id64, err := strconv.ParseUint(val, 10, 64)
	if err != nil || id64 == 0 {
		return 0, err
	}
	return uint(id64), nil
}

func NewProfessorController(professorService service.ProfessorService) *ProfessorController {
	return &ProfessorController{
		professorService: professorService,
	}
}
