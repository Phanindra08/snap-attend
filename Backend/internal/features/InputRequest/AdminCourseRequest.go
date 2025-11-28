package InputRequest

type CreateCourseRequest struct {
	CourseName string `json:"courseName" binding:"required"`
}

type UpdateCourseRequest struct {
	CourseName string `json:"courseName" binding:"required"`
}
