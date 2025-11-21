package service

import (
	"context"
	"errors"
	"strings"

	"github.com/phanindra08/snap-attend/internal/features/dto"
	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidProfile        = errors.New("invalid profile")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrNoFieldsToUpdate      = errors.New("no fields to update")
	ErrProfileNotConfigured  = errors.New("user profile not configured in database")
	ErrOldPasswordIncorrect  = errors.New("old password is incorrect")
	ErrPasswordFieldsInvalid = errors.New("both oldPassword and newPassword are required")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters")
)

type UserService interface {
	RegisterUser(ctx context.Context, signUpDTO dto.SignupDTO) (*models.User, error)
	LoginUser(ctx context.Context, loginDTO dto.LoginDTO) (*models.User, models.UserRoles, error)
	UpdateUser(ctx context.Context, userID uint, updateUserDTO dto.UpdateUserDTO) (*models.User, error)
	FindUsersByName(ctx context.Context, firstName string, lastName string) ([]models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (userService *userService) RegisterUser(ctx context.Context, signUpDTO dto.SignupDTO) (*models.User, error) {
	if signUpDTO.Profile != models.Student && signUpDTO.Profile != models.Professor {
		return nil, ErrInvalidProfile
	}

	country := strings.TrimSpace(signUpDTO.Country)
	if utils.IsEmpty(country) {
		country = utils.DEFAULT_COUNTRY
	}

	// Ensuring UserProfile exists
	profile, err := userService.userRepo.GetUserProfileByRole(ctx, signUpDTO.Profile)
	if err != nil {
		if errors.Is(err, repository.ErrUserProfileNotFound) {
			return nil, ErrProfileNotConfigured
		}
		return nil, err
	}

	// Creating Address
	addr := &models.Address{
		Address1: signUpDTO.Address1,
		Address2: signUpDTO.Address2,
		City:     signUpDTO.City,
		State:    signUpDTO.State,
		ZipCode:  signUpDTO.ZipCode,
		Country:  country,
	}

	if err := userService.userRepo.SaveAddress(ctx, addr); err != nil {
		return nil, err
	}

	// Creating User
	user := &models.User{
		FirstName: signUpDTO.FirstName,
		LastName:  signUpDTO.LastName,
		Email:     signUpDTO.Email,
		Password:  signUpDTO.Password,
		ProfileId: profile.ID,
		AddressId: addr.ID,
	}

	// Hashing password using Argon2
	if err := user.PasswordHashing(); err != nil {
		return nil, err
	}

	// Save user
	if err := userService.userRepo.SaveUser(ctx, user); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (userService *userService) LoginUser(ctx context.Context, loginDTO dto.LoginDTO) (*models.User, models.UserRoles, error) {
	if loginDTO.Profile != models.Student &&
		loginDTO.Profile != models.Professor &&
		loginDTO.Profile != models.Admin {
		return nil, "", ErrInvalidProfile
	}

	user, err := userService.userRepo.FindByEmailAndRole(ctx, loginDTO.Email, loginDTO.Profile)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if !user.ComparePassword(loginDTO.Password) {
		return nil, "", ErrInvalidCredentials
	}

	return user, loginDTO.Profile, nil
}

func (userService *userService) FindUsersByName(ctx context.Context, firstName string, lastName string) ([]models.User, error) {
	users, err := userService.userRepo.FindByName(ctx, firstName, lastName)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userService *userService) UpdateUser(ctx context.Context, userID uint, updateUserDTO dto.UpdateUserDTO) (*models.User, error) {
	user, err := userService.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}

	isUserUpdated := false
	if !utils.IsEmpty(updateUserDTO.FirstName) {
		user.FirstName = strings.TrimSpace(updateUserDTO.FirstName)
		isUserUpdated = true
	}
	if !utils.IsEmpty(updateUserDTO.LastName) {
		user.LastName = strings.TrimSpace(updateUserDTO.LastName)
		isUserUpdated = true
	}

	oldPassword := strings.TrimSpace(updateUserDTO.OldPassword)
	newPassword := strings.TrimSpace(updateUserDTO.NewPassword)

	if !utils.IsEmpty(oldPassword) || !utils.IsEmpty(newPassword) {

		// Checking whether both password fields are provided
		if utils.IsEmpty(oldPassword) || utils.IsEmpty(newPassword) {
			return nil, ErrPasswordFieldsInvalid
		}

		// checking whether old password is correct
		if !user.ComparePassword(oldPassword) {
			return nil, ErrOldPasswordIncorrect
		}

		if len(newPassword) < 8 {
			return nil, ErrPasswordTooShort
		}

		user.Password = newPassword
		if err := user.PasswordHashing(); err != nil {
			return nil, err
		}
		isUserUpdated = true
	}

	isAddressUpdated := false
	addrFieldsPresent := !utils.IsEmpty(updateUserDTO.Address1) || !utils.IsEmpty(updateUserDTO.Address2) ||
		!utils.IsEmpty(updateUserDTO.City) || !utils.IsEmpty(updateUserDTO.State) ||
		!utils.IsEmpty(updateUserDTO.ZipCode) || !utils.IsEmpty(updateUserDTO.Country)

	if addrFieldsPresent {
		addr, err := userService.userRepo.GetAddressByID(ctx, user.AddressId)
		if err != nil {
			return nil, err
		}

		if !utils.IsEmpty(updateUserDTO.Address1) {
			addr.Address1 = strings.TrimSpace(updateUserDTO.Address1)
			isAddressUpdated = true
		}
		if !utils.IsEmpty(updateUserDTO.Address2) {
			addr.Address2 = strings.TrimSpace(updateUserDTO.Address2)
			isAddressUpdated = true
		}
		if !utils.IsEmpty(updateUserDTO.City) {
			addr.City = strings.TrimSpace(updateUserDTO.City)
			isAddressUpdated = true
		}
		if !utils.IsEmpty(updateUserDTO.State) {
			addr.State = strings.TrimSpace(updateUserDTO.State)
			isAddressUpdated = true
		}
		if !utils.IsEmpty(updateUserDTO.ZipCode) {
			addr.ZipCode = strings.TrimSpace(updateUserDTO.ZipCode)
			isAddressUpdated = true
		}
		if !utils.IsEmpty(updateUserDTO.Country) {
			addr.Country = strings.TrimSpace(updateUserDTO.Country)
			isAddressUpdated = true
		}

		if isAddressUpdated {
			if err := userService.userRepo.UpdateAddress(ctx, addr); err != nil {
				return nil, err
			}
			isUserUpdated = true
		}
	}

	if !isUserUpdated {
		return nil, ErrNoFieldsToUpdate
	}

	if err := userService.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
