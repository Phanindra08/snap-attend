package InputRequest

type StudentSubmitAttendanceRequest struct {
	QrId        uint    `json:"qrId" binding:"required"`
	StudentName string  `json:"studentName" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Question    string  `json:"question"`
	Answer      string  `json:"answer"`
}
