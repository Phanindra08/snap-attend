package InputRequest

type CreateRoomRequest struct {
	RoomNumber   string  `json:"roomNumber" binding:"required"`
	BuildingName string  `json:"buildingName" binding:"required"`
	Capacity     int     `json:"capacity" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
}

type UpdateRoomRequest struct {
	RoomNumber   string  `json:"roomNumber" binding:"required"`
	BuildingName string  `json:"buildingName" binding:"required"`
	Capacity     int     `json:"capacity" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
}
