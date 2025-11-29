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
	ErrRoomAlreadyExists      = errors.New("room already exists")
	ErrRoomInvalidCapacity    = errors.New("room capacity must be greater than zero")
	ErrRoomInvalidCoordinates = errors.New("invalid room coordinates")
	ErrRoomFieldsInvalid      = errors.New("one or more room fields are invalid")
)

type RoomService interface {
	CreateRoom(ctx context.Context, roomDTO *dto.RoomDTO) (*models.Room, error)
	GetRoomByID(ctx context.Context, id uint) (*models.Room, error)
	UpdateRoom(ctx context.Context, id uint, roomDTO *dto.RoomDTO) (*models.Room, error)
	DeleteRoom(ctx context.Context, id uint) error
}

type roomService struct {
	roomRepo repository.RoomRepository
}

func (rs *roomService) validateRoomDTO(roomDTO *dto.RoomDTO) error {
	if roomDTO == nil {
		return ErrRoomFieldsInvalid
	}

	capacity := roomDTO.Capacity
	if utils.IsEmpty(roomDTO.RoomNumber) || utils.IsEmpty(roomDTO.BuildingName) || utils.IsEmpty(roomDTO.Address) {
		return ErrRoomFieldsInvalid
	}
	if capacity <= 0 {
		return ErrRoomInvalidCapacity
	}
	if !utils.IsValidLatAndLon(roomDTO.Latitude, roomDTO.Longitude) {
		return ErrRoomInvalidCoordinates
	}
	return nil
}

func (rs *roomService) CreateRoom(ctx context.Context, roomDTO *dto.RoomDTO) (*models.Room, error) {
	// Validating the DTO fields
	if err := rs.validateRoomDTO(roomDTO); err != nil {
		return nil, err
	}

	roomNumber := strings.TrimSpace(roomDTO.RoomNumber)
	buildingName := strings.TrimSpace(roomDTO.BuildingName)

	// Checking for existing room with same number and building
	existing, err := rs.roomRepo.GetRoomByNumberAndBuilding(ctx, roomNumber, buildingName)
	if err == nil && existing != nil {
		return nil, ErrRoomAlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrRoomNotFound) {
		return nil, err
	}

	room := &models.Room{
		RoomNumber:   roomNumber,
		BuildingName: buildingName,
		Capacity:     roomDTO.Capacity,
		Address:      strings.TrimSpace(roomDTO.Address),
		Latitude:     roomDTO.Latitude,
		Longitude:    roomDTO.Longitude,
	}

	if err := rs.roomRepo.CreateRoom(ctx, room); err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") {
			return nil, ErrRoomAlreadyExists
		}
		return nil, err
	}
	return room, nil
}

func (rs *roomService) GetRoomByID(ctx context.Context, id uint) (*models.Room, error) {
	room, err := rs.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rs *roomService) UpdateRoom(ctx context.Context, id uint, roomDTO *dto.RoomDTO) (*models.Room, error) {
	// Fetch existing room details
	room, err := rs.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validating the DTO fields
	if err := rs.validateRoomDTO(roomDTO); err != nil {
		return nil, err
	}

	roomNumber := strings.TrimSpace(roomDTO.RoomNumber)
	buildingName := strings.TrimSpace(roomDTO.BuildingName)

	// Check uniqueness if the room number or building changed
	if room.RoomNumber != roomNumber || room.BuildingName != buildingName {
		existing, err := rs.roomRepo.GetRoomByNumberAndBuilding(ctx, roomNumber, buildingName)
		if err == nil && existing != nil && existing.ID != room.ID {
			return nil, ErrRoomAlreadyExists
		}
		if err != nil && !errors.Is(err, repository.ErrRoomNotFound) {
			return nil, err
		}
	}

	room.RoomNumber = roomNumber
	room.BuildingName = buildingName
	room.Address = strings.TrimSpace(roomDTO.Address)
	room.Capacity = roomDTO.Capacity
	room.Latitude = roomDTO.Latitude
	room.Longitude = roomDTO.Longitude

	if err := rs.roomRepo.UpdateRoom(ctx, room); err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") {
			return nil, ErrRoomAlreadyExists
		}
		return nil, err
	}
	return room, nil
}

func (rs *roomService) DeleteRoom(ctx context.Context, id uint) error {
	// Ensuring the room record exists
	_, err := rs.roomRepo.GetRoomByID(ctx, id)
	if err != nil {
		return err
	}
	return rs.roomRepo.DeleteRoom(ctx, id)
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
	return &roomService{roomRepo: roomRepo}
}
