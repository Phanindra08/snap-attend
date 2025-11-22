package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phanindra08/snap-attend/internal/features/InputRequest"
	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

type AdminController struct {
	userService       service.UserService
	userRepo          repository.UserRepository
	courseService     service.CourseService
	sectionService    service.SectionService
	enrollmentService service.EnrollmentService
}

func (ac *AdminController) getUserEnsuringRole(context *gin.Context, userID uint, expectedRole models.UserRoles) (*models.User, error) {
	ctx := context.Request.Context()

	user, err := ac.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	profile, err := ac.userRepo.GetUserProfileByRole(ctx, expectedRole)
	if err != nil {
		return nil, err
	}

	if user.ProfileId != profile.ID {
		return nil, service.ErrInvalidProfile
	}

	return user, nil
}

func (ac *AdminController) sanitizeAdminCreateUserRequest(createUserRequest *InputRequest.AdminCreateUserRequest) {
	createUserRequest.Email = utils.TrimAndConvertToLowerCase(createUserRequest.Email)
	createUserRequest.FirstName = strings.TrimSpace(createUserRequest.FirstName)
	createUserRequest.LastName = strings.TrimSpace(createUserRequest.LastName)
	createUserRequest.Password = strings.TrimSpace(createUserRequest.Password)
	createUserRequest.Address1 = strings.TrimSpace(createUserRequest.Address1)
	createUserRequest.Address2 = strings.TrimSpace(createUserRequest.Address2)
	createUserRequest.City = strings.TrimSpace(createUserRequest.City)
	createUserRequest.State = strings.TrimSpace(createUserRequest.State)
	createUserRequest.ZipCode = strings.TrimSpace(createUserRequest.ZipCode)
	createUserRequest.Country = strings.TrimSpace(createUserRequest.Country)
}

func (ac *AdminController) sanitizeAdminUpdateUserRequest(updateUserRequest *InputRequest.AdminUpdateUserRequest) {
	updateUserRequest.FirstName = strings.TrimSpace(updateUserRequest.FirstName)
	updateUserRequest.LastName = strings.TrimSpace(updateUserRequest.LastName)
	updateUserRequest.Address1 = strings.TrimSpace(updateUserRequest.Address1)
	updateUserRequest.Address2 = strings.TrimSpace(updateUserRequest.Address2)
	updateUserRequest.City = strings.TrimSpace(updateUserRequest.City)
	updateUserRequest.State = strings.TrimSpace(updateUserRequest.State)
	updateUserRequest.ZipCode = strings.TrimSpace(updateUserRequest.ZipCode)
	updateUserRequest.Country = strings.TrimSpace(updateUserRequest.Country)
}

func (ac *AdminController) CreateStudent(ctx *gin.Context) {
	ac.createUserWithRole(ctx, models.Student, "Student created successfully")
}

func (ac *AdminController) CreateProfessor(ctx *gin.Context) {
	ac.createUserWithRole(ctx, models.Professor, "Professor created successfully")
}

func (ac *AdminController) CreateAdmin(ctx *gin.Context) {
	ac.createUserWithRole(ctx, models.Admin, "Admin created successfully")
}

func (ac *AdminController) createUserWithRole(ctx *gin.Context, role models.UserRoles, successMessage string) {
	var createUserRequest InputRequest.AdminCreateUserRequest
	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ac.sanitizeAdminCreateUserRequest(&createUserRequest)

	input := dto.SignupDTO{
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		Email:     createUserRequest.Email,
		Password:  createUserRequest.Password,
		Profile:   role,

		Address1: createUserRequest.Address1,
		Address2: createUserRequest.Address2,
		City:     createUserRequest.City,
		State:    createUserRequest.State,
		ZipCode:  createUserRequest.ZipCode,
		Country:  createUserRequest.Country,
	}

	user, err := ac.userService.RegisterUser(ctx.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile"})
			return
		}
		if errors.Is(err, service.ErrProfileNotConfigured) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User profile not configured. Please contact administrator."})
			return
		}

		log.Printf("error while creating user by Admin: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusCreated, gin.H{
		"message": successMessage,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"profile":   role,
		},
	})
}

func (ac *AdminController) GetStudentByID(ctx *gin.Context) {
	ac.getUserByIDAndRole(ctx, models.Student, "Student")
}

func (ac *AdminController) GetProfessorByID(ctx *gin.Context) {
	ac.getUserByIDAndRole(ctx, models.Professor, "Professor")
}

func (ac *AdminController) getUserByIDAndRole(ctx *gin.Context, role models.UserRoles, label string) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	user, err := ac.getUserEnsuringRole(ctx, id, role)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": label + " not found"})
			return
		}
		if errors.Is(err, repository.ErrUserProfileNotFound) || errors.Is(err, service.ErrProfileNotConfigured) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User profile not configured. Please contact administrator."})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is not a " + label})
			return
		}
		log.Printf("Admin get %s error: %v", label, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve " + strings.ToLower(label)})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"profile":   role,
	})
}

func (ac *AdminController) UpdateStudent(ctx *gin.Context) {
	ac.updateUserWithRole(ctx, models.Student, "Student")
}

func (ac *AdminController) UpdateProfessor(ctx *gin.Context) {
	ac.updateUserWithRole(ctx, models.Professor, "Professor")
}

func (ac *AdminController) updateUserWithRole(ctx *gin.Context, role models.UserRoles, label string) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Ensure the user exists and has the correct role
	_, err = ac.getUserEnsuringRole(ctx, id, role)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": label + " not found"})
			return
		}
		if errors.Is(err, repository.ErrUserProfileNotFound) || errors.Is(err, service.ErrProfileNotConfigured) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User profile not configured. Please contact administrator."})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is not a " + label})
			return
		}
		log.Printf("Admin get %s for update error: %v", label, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve " + strings.ToLower(label)})
		return
	}

	var req InputRequest.AdminUpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ac.sanitizeAdminUpdateUserRequest(&req)

	updateDTO := dto.UpdateUserDTO{
		FirstName: req.FirstName,
		LastName:  req.LastName,

		Address1: req.Address1,
		Address2: req.Address2,
		City:     req.City,
		State:    req.State,
		ZipCode:  req.ZipCode,
		Country:  req.Country,
	}

	updatedUser, err := ac.userService.UpdateUser(ctx.Request.Context(), id, updateDTO)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoFieldsToUpdate):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
			return
		case errors.Is(err, repository.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": label + " not found"})
			return
		default:
			log.Printf("Admin update %s error: %v", label, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update " + strings.ToLower(label)})
			return
		}
	}

	updatedUser.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"message":   label + " updated successfully",
		"id":        updatedUser.ID,
		"email":     updatedUser.Email,
		"firstName": updatedUser.FirstName,
		"lastName":  updatedUser.LastName,
		"profile":   role,
	})
}

func (ac *AdminController) CreateCourse(ctx *gin.Context) {
	var createCourseRequest InputRequest.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&createCourseRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(createCourseRequest.CourseName)

	course, err := ac.courseService.CreateCourse(ctx.Request.Context(), name)
	if err != nil {
		if errors.Is(err, service.ErrCourseAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course already exists"})
			return
		}
		log.Printf("Admin create course error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course": gin.H{
			"id":         course.ID,
			"courseName": course.CourseName,
		},
	})
}

func (ac *AdminController) GetCourseByID(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	course, err := ac.courseService.GetCourseByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		log.Printf("Error retrieving course error by admin: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve course"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         course.ID,
		"courseName": course.CourseName,
	})
}

func (ac *AdminController) UpdateCourse(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var updateCourseRequest InputRequest.UpdateCourseRequest
	if err := ctx.ShouldBindJSON(&updateCourseRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(updateCourseRequest.CourseName)

	course, err := ac.courseService.UpdateCourse(ctx.Request.Context(), id, name)
	if err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		if errors.Is(err, service.ErrCourseAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course already exists"})
			return
		}
		log.Printf("Error while updating course by admin: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Course updated successfully",
		"id":         course.ID,
		"courseName": course.CourseName,
	})
}

func (ac *AdminController) DeleteCourse(ctx *gin.Context) {
	id, err := parseIDParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	if err := ac.courseService.DeleteCourse(ctx.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrCourseNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		log.Printf("Error while deleting the course by admin: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course deleted successfully",
	})
}

func (ac *AdminController) EnrollStudentInSection(ctx *gin.Context) {
	var adminEnrollStudentRequest InputRequest.AdminEnrollStudentRequest
	if err := ctx.ShouldBindJSON(&adminEnrollStudentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ac.getUserEnsuringRole(ctx, adminEnrollStudentRequest.StudentId, models.Student)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		if errors.Is(err, repository.ErrUserProfileNotFound) || errors.Is(err, service.ErrProfileNotConfigured) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User profile not configured. Please contact administrator."})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is not a Student"})
			return
		}
		log.Printf("Error verifying student before enrollment: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify student"})
		return
	}

	enrollment, err := ac.enrollmentService.EnrollStudentInSection(ctx.Request.Context(), adminEnrollStudentRequest.StudentId, adminEnrollStudentRequest.SectionId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrSectionNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Course section not found"})
			return
		case errors.Is(err, service.ErrMaxCoursesPerSemester):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student has reached maximum number of courses for this semester"})
			return
		case errors.Is(err, service.ErrAlreadyEnrolledInCourse):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student is already enrolled in this course for the semester"})
			return
		case errors.Is(err, service.ErrStudentScheduleConflict):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Student has a schedule conflict with another course"})
			return
		case errors.Is(err, service.ErrSectionFull):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Course section is full"})
			return
		default:
			log.Printf("Admin enroll student error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll student in course section"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Student enrolled successfully",
		"enrollment": gin.H{
			"id":        enrollment.ID,
			"studentId": enrollment.StudentId,
			"sectionId": enrollment.SectionId,
		},
	})
}

func (ac *AdminController) AssignProfessorToSection(ctx *gin.Context) {
	var assignProfessorRequest InputRequest.AdminAssignProfessorRequest
	if err := ctx.ShouldBindJSON(&assignProfessorRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := ac.getUserEnsuringRole(ctx, assignProfessorRequest.ProfessorId, models.Professor)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Professor not found"})
			return
		}
		if errors.Is(err, repository.ErrUserProfileNotFound) || errors.Is(err, service.ErrProfileNotConfigured) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User profile not configured. Please contact administrator."})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User is not a Professor"})
			return
		}
		log.Printf("Error verifying professor before assignment: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify professor"})
		return
	}

	section, err := ac.sectionService.AssignProfessorToSection(ctx.Request.Context(), assignProfessorRequest.SectionId, assignProfessorRequest.ProfessorId)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrSectionNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Course section not found"})
			return
		case errors.Is(err, service.ErrProfessorScheduleConflict):
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Professor has a schedule conflict with another section"})
			return
		default:
			log.Printf("Admin assign professor error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign professor to course section"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Professor assigned successfully",
		"section": gin.H{
			"id":            section.ID,
			"courseId":      section.CourseId,
			"professorId":   section.ProfessorId,
			"sectionNumber": section.SectionNumber,
		},
	})
}

func (ac *AdminController) Search(ctx *gin.Context) {
	courseName := strings.TrimSpace(ctx.Query("courseName"))
	studentName := strings.TrimSpace(ctx.Query("studentName"))
	professorName := strings.TrimSpace(ctx.Query("professorName"))

	if utils.IsEmpty(courseName) && utils.IsEmpty(studentName) && utils.IsEmpty(professorName) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "At least one search parameter must be provided"})
		return
	}

	response := gin.H{}

	if !utils.IsEmpty(courseName) {
		courses, err := ac.courseService.SearchCoursesByName(ctx.Request.Context(), courseName)
		if err != nil {
			log.Printf("Error while searching courses by admin: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search courses"})
			return
		}

		var courseResults []gin.H
		for _, course := range courses {
			courseResults = append(courseResults, gin.H{
				"id":         course.ID,
				"courseName": course.CourseName,
			})
		}
		response["courses"] = courseResults
	}

	if !utils.IsEmpty(studentName) {
		students, err := ac.userRepo.SearchUsersByNameAndRole(ctx.Request.Context(), studentName, models.Student)
		if err != nil {
			log.Printf("Error while searching students by admin: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search students"})
			return
		}

		var studentResults []gin.H
		for _, student := range students {
			studentResults = append(studentResults, gin.H{
				"id":        student.ID,
				"email":     student.Email,
				"firstName": student.FirstName,
				"lastName":  student.LastName,
				"profile":   models.Student,
			})
		}
		response["students"] = studentResults
	}

	if !utils.IsEmpty(professorName) {
		professors, err := ac.userRepo.SearchUsersByNameAndRole(ctx.Request.Context(), professorName, models.Professor)
		if err != nil {
			log.Printf("Error while searching professors by admin: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search professors"})
			return
		}

		var professorResults []gin.H
		for _, professor := range professors {
			professorResults = append(professorResults, gin.H{
				"id":        professor.ID,
				"email":     professor.Email,
				"firstName": professor.FirstName,
				"lastName":  professor.LastName,
				"profile":   models.Professor,
			})
		}
		response["professors"] = professorResults
	}

	ctx.JSON(http.StatusOK, response)
}

func parseIDParam(ctx *gin.Context) (uint, error) {
	idParam := ctx.Param("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

func NewAdminController(
	userService service.UserService,
	userRepo repository.UserRepository,
	courseService service.CourseService,
	sectionService service.SectionService,
	enrollmentService service.EnrollmentService,
) *AdminController {
	return &AdminController{
		userService:       userService,
		userRepo:          userRepo,
		courseService:     courseService,
		sectionService:    sectionService,
		enrollmentService: enrollmentService,
	}
}
