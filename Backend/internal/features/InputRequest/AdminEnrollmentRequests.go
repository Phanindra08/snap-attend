package InputRequest

type AdminEnrollStudentRequest struct {
	StudentId uint `json:"studentId" binding:"required"`
	SectionId uint `json:"sectionId" binding:"required"`
}

type AdminAssignProfessorRequest struct {
	SectionId   uint `json:"sectionId" binding:"required"`
	ProfessorId uint `json:"professorId" binding:"required"`
}
