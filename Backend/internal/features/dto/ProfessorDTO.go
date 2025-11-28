package dto

import (
	"time"

	"github.com/phanindra08/snap-attend/internal/shared/models"
)

type DailyAttendanceRow struct {
	Student    models.User
	Present    bool
	AttendedAt *time.Time
}

type CourseSearchResult struct {
	CourseID      uint
	CourseName    string
	SectionID     uint
	SectionNumber string
}

type StudentSearchResult struct {
	StudentID uint
	FirstName string
	LastName  string
	Email     string
}

type ProfessorSearchResult struct {
	Courses  []CourseSearchResult
	Students []StudentSearchResult
}

type SectionAttendanceOverviewRow struct {
	StudentID       uint
	FirstName       string
	LastName        string
	Email           string
	AttendedClasses int64
	TotalClasses    int64
	AttendancePct   float64
}

type SectionAttendanceOverview struct {
	SectionID     uint
	TotalSessions int64
	Students      []SectionAttendanceOverviewRow
}
