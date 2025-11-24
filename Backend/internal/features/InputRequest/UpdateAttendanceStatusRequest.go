package InputRequest

type UpdateAttendanceStatusRequest struct {
	Attended *bool `json:"attended" binding:"required"`
}
