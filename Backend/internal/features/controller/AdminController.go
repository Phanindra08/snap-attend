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
	userService   service.UserService
	userRepo      repository.UserRepository
	courseService service.CourseService
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
) *AdminController {
	return &AdminController{
		userService:   userService,
		userRepo:      userRepo,
		courseService: courseService,
	}
}
