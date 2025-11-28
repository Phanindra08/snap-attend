package InputRequest

type StudentSubmitAttendanceRequest struct {
	QrId      uint    `json:"qrId" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Question  string  `json:"question"`
	Answer    string  `json:"answer"`
}
