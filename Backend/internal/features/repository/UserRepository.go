package repository

import (
	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func (userRepo *userRepository) SaveUser(user *models.User) error {
	return userRepo.db.Create(user).Error
}

func (userRepo *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := userRepo.db.First(&user, id).Error
	return &user, err
}

func (userRepo *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := userRepo.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (userRepo *userRepository) UpdateUser(user *models.User) error {
	return userRepo.db.Save(user).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
