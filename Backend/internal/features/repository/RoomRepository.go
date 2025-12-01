package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound = errors.New("room not found")
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	GetRoomByID(ctx context.Context, id uint) (*models.Room, error)
	GetRoomByNumberAndBuilding(ctx context.Context, roomNumber string, buildingName string) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id uint) error
}

type roomRepository struct {
	db *gorm.DB
}

func (roomRepo *roomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	if err := roomRepo.db.WithContext(ctx).Create(room).Error; err != nil {
		return fmt.Errorf("can't create the room record: %w", err)
	}
	return nil
}

func (roomRepo *roomRepository) GetRoomByID(ctx context.Context, id uint) (*models.Room, error) {
	var roomDetails models.Room
	if err := roomRepo.db.WithContext(ctx).First(&roomDetails, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("can't find the room details: %w", err)
	}
	return &roomDetails, nil
}

func (roomRepo *roomRepository) GetRoomByNumberAndBuilding(ctx context.Context, roomNumber string, buildingName string) (*models.Room, error) {
	var roomDetails models.Room
	if err := roomRepo.db.WithContext(ctx).
		Where("room_number = ? AND building_name = ?", roomNumber, buildingName).
		First(&roomDetails).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, fmt.Errorf("can't find the room details: %w", err)
	}
	return &roomDetails, nil
}

func (roomRepo *roomRepository) UpdateRoom(ctx context.Context, room *models.Room) error {
	if err := roomRepo.db.WithContext(ctx).Save(room).Error; err != nil {
		return fmt.Errorf("can't update the room details: %w", err)
	}
	return nil
}

func (roomRepo *roomRepository) DeleteRoom(ctx context.Context, id uint) error {
	if err := roomRepo.db.WithContext(ctx).Delete(&models.Room{}, id).Error; err != nil {
		return fmt.Errorf("can't delete the room details: %w", err)
	}
	return nil
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}
