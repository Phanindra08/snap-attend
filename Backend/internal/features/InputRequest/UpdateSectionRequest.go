package InputRequest

type UpdateSectionRequest struct {
	SectionNumber string   `json:"sectionNumber"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	RoomId        *uint    `json:"roomId"`
}
