package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phanindra08/snap-attend/internal/features/InputRequest"
	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/middleware"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

type UserController struct {
	userService service.UserService
}

func (uc *UserController) Signup(ctx *gin.Context) {
	var signupRequest InputRequest.SignupRequest
	if err := ctx.ShouldBindJSON(&signupRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trimming whitespaces at the beginning and end of input fields
	signupRequest.Email = utils.TrimAndConvertToLowerCase(signupRequest.Email)
	signupRequest.FirstName = strings.TrimSpace(signupRequest.FirstName)
	signupRequest.LastName = strings.TrimSpace(signupRequest.LastName)
	signupRequest.Password = strings.TrimSpace(signupRequest.Password)
	signupRequest.Address1 = strings.TrimSpace(signupRequest.Address1)
	signupRequest.Address2 = strings.TrimSpace(signupRequest.Address2)
	signupRequest.City = strings.TrimSpace(signupRequest.City)
	signupRequest.State = strings.TrimSpace(signupRequest.State)
	signupRequest.ZipCode = strings.TrimSpace(signupRequest.ZipCode)
	signupRequest.Country = strings.TrimSpace(signupRequest.Country)

	if signupRequest.Profile != models.Student && signupRequest.Profile != models.Professor {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Profile must be Student or Professor"})
		return
	}

	input := dto.SignupDTO{
		FirstName: signupRequest.FirstName,
		LastName:  signupRequest.LastName,
		Email:     signupRequest.Email,
		Password:  signupRequest.Password,
		Profile:   signupRequest.Profile,

		Address1: signupRequest.Address1,
		Address2: signupRequest.Address2,
		City:     signupRequest.City,
		State:    signupRequest.State,
		ZipCode:  signupRequest.ZipCode,
		Country:  signupRequest.Country,
	}

	user, err := uc.userService.RegisterUser(ctx.Request.Context(), input)
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

		log.Printf("Signup error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"profile":   signupRequest.Profile,
		},
	})
}

func (uc *UserController) Login(c *gin.Context) {
	var loginRequest InputRequest.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginRequest.Email = utils.TrimAndConvertToLowerCase(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)

	user, role, err := uc.userService.LoginUser(c.Request.Context(), dto.LoginDTO{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
		Profile:  loginRequest.Profile,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		if errors.Is(err, service.ErrInvalidProfile) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile"})
			return
		}
		log.Printf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, role)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"profile":   role,
		},
	})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID") // set by JWTAuthMiddleware

	var updateProfileRequest InputRequest.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trimming whitespaces at the beginning and end of input fields
	updateProfileRequest.FirstName = strings.TrimSpace(updateProfileRequest.FirstName)
	updateProfileRequest.LastName = strings.TrimSpace(updateProfileRequest.LastName)
	updateProfileRequest.OldPassword = strings.TrimSpace(updateProfileRequest.OldPassword)
	updateProfileRequest.NewPassword = strings.TrimSpace(updateProfileRequest.NewPassword)
	updateProfileRequest.Address1 = strings.TrimSpace(updateProfileRequest.Address1)
	updateProfileRequest.Address2 = strings.TrimSpace(updateProfileRequest.Address2)
	updateProfileRequest.City = strings.TrimSpace(updateProfileRequest.City)
	updateProfileRequest.State = strings.TrimSpace(updateProfileRequest.State)
	updateProfileRequest.ZipCode = strings.TrimSpace(updateProfileRequest.ZipCode)
	updateProfileRequest.Country = strings.TrimSpace(updateProfileRequest.Country)

	user, err := uc.userService.UpdateUser(
		c.Request.Context(),
		userID,
		dto.UpdateUserDTO{
			FirstName:   updateProfileRequest.FirstName,
			LastName:    updateProfileRequest.LastName,
			OldPassword: updateProfileRequest.OldPassword,
			NewPassword: updateProfileRequest.NewPassword,
			Address1:    updateProfileRequest.Address1,
			Address2:    updateProfileRequest.Address2,
			City:        updateProfileRequest.City,
			State:       updateProfileRequest.State,
			ZipCode:     updateProfileRequest.ZipCode,
			Country:     updateProfileRequest.Country,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoFieldsToUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
			return
		case errors.Is(err, service.ErrPasswordFieldsInvalid),
			errors.Is(err, service.ErrPasswordTooShort),
			errors.Is(err, service.ErrOldPasswordIncorrect):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			log.Printf("UpdateProfile error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message":   "Profile updated successfully",
		"id":        user.ID,
		"email":     user.Email,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
	})
}

func (uc *UserController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully. Please remove token on client side.",
	})
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}
