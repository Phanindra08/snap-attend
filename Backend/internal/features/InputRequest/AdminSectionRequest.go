package InputRequest

type CreateSectionRequest struct {
	CourseId      uint    `json:"courseId" binding:"required"`
	SectionNumber string  `json:"sectionNumber" binding:"required"`
	SemesterId    uint    `json:"semesterId" binding:"required"`
	RoomId        uint    `json:"roomId" binding:"required"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
}
