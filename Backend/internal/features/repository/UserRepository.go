package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserProfileNotFound = errors.New("user profile not found")
)

type UserRepository interface {
	GetUserProfileByRole(ctx context.Context, role models.UserRoles) (*models.UserProfile, error)
	SaveAddress(ctx context.Context, address *models.Address) error
	SaveUser(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmailAndRole(ctx context.Context, email string, role models.UserRoles) (*models.User, error)
	FindByName(ctx context.Context, firstName string, lastName string) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

// ----- DB Operations for UserProfile -----

func (userRepo *userRepository) GetUserProfileByRole(ctx context.Context, role models.UserRoles) (*models.UserProfile, error) {
	var profile models.UserProfile
	err := userRepo.db.WithContext(ctx).Where("role = ?", role).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserProfileNotFound
		}
		return nil, fmt.Errorf("can't find the user role: %w", err)
	}
	return &profile, nil
}

// ----- DB Operations for Address -----

func (userRepo *userRepository) SaveAddress(ctx context.Context, address *models.Address) error {
	err := userRepo.db.WithContext(ctx).Create(address).Error
	if err != nil {
		return fmt.Errorf("can't save address of the user: %w", err)
	}
	return nil
}

// ----- DB Operations for User -----

func (userRepo *userRepository) SaveUser(ctx context.Context, user *models.User) error {
	err := userRepo.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return fmt.Errorf("can't save user: %w", err)
	}
	return nil
}

func (userRepo *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := userRepo.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("can't find the user: %w", err)
	}
	return &user, nil
}

func (userRepo *userRepository) FindByEmailAndRole(ctx context.Context, email string, role models.UserRoles) (*models.User, error) {
	var user models.User
	err := userRepo.db.WithContext(ctx).
		Joins("UserProfile").
		Where("users.email = ? AND user_profiles.role = ?", email, role).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("can't find the user: %w", err)
	}
	return &user, nil
}

func (userRepo *userRepository) FindByName(ctx context.Context, firstName string, lastName string) ([]models.User, error) {
	var users []models.User
	err := userRepo.db.WithContext(ctx).Where("first_name = ? AND last_name = ?", firstName, lastName).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("can't find the user: %w", err)
	}
	return users, nil
}

func (userRepo *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	err := userRepo.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return fmt.Errorf("can't update the user: %w", err)
	}
	return nil
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
