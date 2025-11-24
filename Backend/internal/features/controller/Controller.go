package controller

import (
	"sync"

	"github.com/phanindra08/snap-attend/internal/features/repository"
	"github.com/phanindra08/snap-attend/internal/features/service"
	"github.com/phanindra08/snap-attend/internal/shared/database"
)

type Controller struct {
	User      *UserController
	Admin     *AdminController
	Student   *StudentController
	Professor *ProfessorController
}

var (
	controllersOnce    sync.Once
	controllerInstance *Controller
)

// GetControllers returns a singleton instance holding all controller pointers.
func GetControllers() *Controller {
	controllersOnce.Do(func() {
		db := database.GetDB()

		// Create all the repositories
		userRepo := repository.NewUserRepository(db)
		courseRepo := repository.NewCourseRepository(db)
		sectionRepo := repository.NewSectionRepository(db)
		enrollmentRepo := repository.NewEnrollmentRepository(db)
		attendanceRepo := repository.NewAttendanceRepository(db)

		// Create all the services
		userService := service.NewUserService(userRepo)
		courseService := service.NewCourseService(courseRepo)
		sectionService := service.NewSectionService(sectionRepo)
		enrollmentService := service.NewEnrollmentService(enrollmentRepo, sectionRepo)
		attendanceService := service.NewAttendanceService(attendanceRepo, enrollmentRepo)

		// Create all the controllers
		controllerInstance = &Controller{
			User: NewUserController(userService),
			Admin: NewAdminController(userService, userRepo, courseService,
				sectionService, enrollmentService),
			Student: NewStudentController(attendanceService),
			Professor: NewProfessorController(attendanceService, sectionRepo, enrollmentRepo,
				userRepo, attendanceRepo),
		}
	})
	return controllerInstance
}
