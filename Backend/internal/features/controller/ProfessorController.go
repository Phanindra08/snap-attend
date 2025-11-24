package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phanindra08/snap-attend/internal/features/InputRequest"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

type ProfessorController struct {
	attendanceService service.AttendanceService
	sectionRepo       repository.SectionRepository
	enrollmentRepo    repository.EnrollmentRepository
	userRepo          repository.UserRepository
	attendanceRepo    repository.AttendanceRepository
}

func (pc *ProfessorController) GenerateAttendanceQr(ctx *gin.Context) {
	professorID := ctx.GetUint("userID")

	sectionID, err := parseUintParam(ctx, "sectionId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sectionId"})
		return
	}

	section, err := pc.sectionRepo.GetSectionByID(ctx.Request.Context(), sectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	if section.ProfessorId != professorID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
		return
	}

	now := time.Now().UTC()
	if !utils.IsClassInSessionNow(now, section.SectionSchedules) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Class is not currently in the scheduled time."})
		return
	}

	// Deactivating any existing active QR's for this section
	if err := pc.attendanceRepo.DeactivateActiveQrsForSection(ctx.Request.Context(), sectionID); err != nil {
		log.Printf("Failed to deactivate previous QRs: %v", err)
	}

	generatedAt := now
	expiresAt := generatedAt.Add(time.Minute * utils.QR_TTL_MINUTES)

	qr := &section.AttendanceQrs
	_ = qr

	newQr := &repository.AttendanceQrModelForCreate{
		SectionId:   sectionID,
		GeneratedAt: generatedAt,
		ExpiresAt:   expiresAt,
	}

	createdQr, err := pc.attendanceRepo.CreateAttendanceQr(ctx.Request.Context(), newQr)
	if err != nil {
		log.Printf("error while generating QR by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate attendance QR"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Attendance QR generated successfully",
		"qr": gin.H{
			"id":          createdQr.ID,
			"sectionId":   createdQr.SectionId,
			"generatedAt": createdQr.GeneratedAt,
			"expiresAt":   createdQr.ExpiresAt,
			"qrHint":      createdQr.QrLink,
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

	section, err := pc.sectionRepo.GetSectionByID(ctx.Request.Context(), sectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}
	if section.ProfessorId != professorID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
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

	fromDate := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	toDate := fromDate.Add(24*time.Hour - time.Nanosecond)

	records, err := pc.attendanceRepo.GetStudentAttendanceBySection(ctx.Request.Context(), sectionID, &fromDate, &toDate, nil)
	if err != nil {
		log.Printf("Error while fetching daily attendance report by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance records"})
		return
	}

	// retrieving all the enrolled students (including those absent that day)
	enrollments, err := pc.enrollmentRepo.GetEnrollmentsBySection(ctx.Request.Context(), sectionID)
	if err != nil {
		log.Printf("Error while fetching daily report based on enrollments by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}

	attendanceByStudent := make(map[uint][]repository.StudentAttendanceWithUser)
	for _, record := range records {
		attendanceByStudent[record.StudentId] = append(attendanceByStudent[record.StudentId], record)
	}

	var result []gin.H
	for _, enrollment := range enrollments {
		student := enrollment.Student
		attendanceRecords := attendanceByStudent[student.ID]

		present := false
		var attendedAt *time.Time

		for _, record := range attendanceRecords {
			if record.Attended {
				present = true
				timeAt := record.AttendedAt
				attendedAt = &timeAt
				break
			}
		}

		result = append(result, gin.H{
			"studentId":  student.ID,
			"firstName":  student.FirstName,
			"lastName":   student.LastName,
			"email":      student.Email,
			"present":    present,
			"attendedAt": attendedAt,
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

	section, err := pc.sectionRepo.GetSectionByID(ctx.Request.Context(), sectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}
	if section.ProfessorId != professorID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
		return
	}

	var (
		from      *time.Time
		to        *time.Time
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
		from = &formattedStartDate
	}

	endDateInStringFormat := strings.TrimSpace(ctx.Query("endDate"))
	if !utils.IsEmpty(endDateInStringFormat) {
		endDate, err := time.Parse("2006-01-02", endDateInStringFormat)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endDate format. Use YYYY-MM-DD"})
			return
		}
		formattedEndDate := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, time.UTC)
		to = &formattedEndDate
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

	records, err := pc.attendanceRepo.GetStudentAttendanceBySection(ctx.Request.Context(), sectionID, from, to, studentId)
	if err != nil {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendanceId"})
		return
	}

	record, err := pc.attendanceRepo.GetAttendanceByIDWithSection(ctx.Request.Context(), attendanceID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Attendance record is not found"})
		return
	}

	section := record.AttendanceQr.CourseSection
	if section.ProfessorId != professorID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
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

	record.Attended = *attendanceStatusRequest.Attended

	if err := pc.attendanceRepo.UpdateStudentAttendance(ctx.Request.Context(), record); err != nil {
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

	if utils.IsEmpty(courseName) || utils.IsEmpty(studentName) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter must be provided"})
		return
	}

	response := gin.H{}

	if !utils.IsEmpty(courseName) {
		sections, err := pc.sectionRepo.SearchSectionsByProfessorAndCourseName(ctx.Request.Context(), professorID, courseName)
		if err != nil {
			log.Printf("Error while searching courses by Professor: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search courses"})
			return
		}

		var courseResults []gin.H
		for _, section := range sections {
			courseResults = append(courseResults, gin.H{
				"courseId":      section.Course.ID,
				"courseName":    section.Course.CourseName,
				"sectionId":     section.ID,
				"sectionNumber": section.SectionNumber,
			})
		}
		response["courses"] = courseResults
	}

	if !utils.IsEmpty(studentName) {
		sections, err := pc.sectionRepo.GetSectionsByProfessor(ctx.Request.Context(), professorID)
		if err != nil {
			log.Printf("Error while searching students sections by Professor: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search students"})
			return
		}

		studentName := utils.TrimAndConvertToLowerCase(studentName)
		seen := make(map[uint]bool)
		var students []gin.H

		for _, section := range sections {
			enrollments, err := pc.enrollmentRepo.GetEnrollmentsBySection(ctx.Request.Context(), section.ID)
			if err != nil {
				log.Printf("Error while searching students enrollments by Professor: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search students"})
				return
			}

			for _, enrollment := range enrollments {
				student := enrollment.Student
				if seen[student.ID] {
					continue
				}

				fullName := utils.TrimAndConvertToLowerCase(student.FirstName + " " + student.LastName)
				if strings.Contains(fullName, studentName) {
					seen[student.ID] = true
					students = append(students, gin.H{
						"studentId": student.ID,
						"firstName": student.FirstName,
						"lastName":  student.LastName,
						"email":     student.Email,
					})
				}
			}
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

	section, err := pc.sectionRepo.GetSectionByID(ctx.Request.Context(), sectionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}
	if section.ProfessorId != professorID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This section is not offered by you. Please select the course which you are teaching."})
		return
	}

	totalSessions, err := pc.attendanceRepo.CountAttendanceQrBySectionUntilNow(ctx.Request.Context(), sectionID)
	if err != nil {
		log.Printf("Error while retrieving the total attendance overview by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count sessions"})
		return
	}

	enrollments, err := pc.enrollmentRepo.GetEnrollmentsBySection(ctx.Request.Context(), sectionID)
	if err != nil {
		log.Printf("Error while retrieving enrollments by Professor: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}

	var result []gin.H
	for _, enrollment := range enrollments {
		student := enrollment.Student

		attendedCount, err := pc.attendanceRepo.CountStudentAttendanceBySection(ctx.Request.Context(), student.ID, sectionID)
		if err != nil {
			log.Printf("Error while retrieving count of student attendance by Professor: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to compute attendance"})
			return
		}

		var pct float64
		if totalSessions > 0 {
			pct = (float64(attendedCount) / float64(totalSessions)) * 100.0
		}

		result = append(result, gin.H{
			"studentId":       student.ID,
			"firstName":       student.FirstName,
			"lastName":        student.LastName,
			"email":           student.Email,
			"attendedClasses": attendedCount,
			"totalClasses":    totalSessions,
			"attendancePct":   pct,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"sectionId":          sectionID,
		"totalSessions":      totalSessions,
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

func NewProfessorController(
	attendanceService service.AttendanceService,
	sectionRepo repository.SectionRepository,
	enrollmentRepo repository.EnrollmentRepository,
	userRepo repository.UserRepository,
	attendanceRepo repository.AttendanceRepository,
) *ProfessorController {
	return &ProfessorController{
		attendanceService: attendanceService,
		sectionRepo:       sectionRepo,
		enrollmentRepo:    enrollmentRepo,
		userRepo:          userRepo,
		attendanceRepo:    attendanceRepo,
	}
}
