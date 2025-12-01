package dto

type StudentSubmitAttendanceDTO struct {
	QrId      uint
	Latitude  float64
	Longitude float64
	Question    string
	Answer      string
	StudentName string
}
