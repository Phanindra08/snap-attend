package database

import (
	"fmt"
	"log"
	"time"

	"github.com/phanindra08/snap-attend/internal/shared/models"
	"github.com/phanindra08/snap-attend/internal/shared/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// initDB Function - Initializes the database connection
func initDB() {
	var err error
	
	// Use SQLite for local development to avoid Postgres dependency
	// dsn := config.GetConfig().GetConnectionStringForDB()

	// Configuring the GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(utils.DB_LOG_LEVEL),
		NowFunc: func() time.Time {
			return time.Now().UTC() // Use UTC for all timestamps
		},
	}

	// Trying to establish DB connection with a retry logic
	for i := 0; i < utils.MAX_DB_RETRIES; i++ {
		db, err = gorm.Open(sqlite.Open("snapattend.db"), gormConfig)
		if err == nil {
			break
		}

		log.Printf("Attempt %d: Failed to connect to the database. Retrying the connection in 5 seconds...", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to the database after %d attempts: %v", utils.MAX_DB_RETRIES, err)
	}

	// Configuring the DB connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve sql connection pool: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(utils.MAX_IDLE_CONNECTIONS)
	sqlDB.SetMaxOpenConns(utils.MAX_OPEN_CONNECTIONS)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(utils.CONNECTIONS_MAX_IDLE_TIME)

	log.Println("Database connection established successfully")
}

// autoMigrate runs database migrations
func autoMigrate() {
	// First we should create independent tables then dependent tables
	err := db.AutoMigrate(&models.Address{})
	if err != nil {
		log.Fatalf("Failed to migrate the table Address: %v", err)
	}

	err = db.AutoMigrate(&models.UserProfile{})
	if err != nil {
		log.Fatalf("Failed to migrate the table UserProfile: %v", err)
	}

	err = db.AutoMigrate(&models.Course{})
	if err != nil {
		log.Fatalf("Failed to migrate the table Course: %v", err)
	}

	err = db.AutoMigrate(&models.Room{})
	if err != nil {
		log.Fatalf("Failed to migrate the table Room: %v", err)
	}

	err = db.AutoMigrate(&models.Semester{})
	if err != nil {
		log.Fatalf("Failed to migrate the table Semester: %v", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate the table users: %v", err)
	}

	err = db.AutoMigrate(&models.CourseSection{})
	if err != nil {
		log.Fatalf("Failed to migrate the table CourseSection: %v", err)
	}

	err = db.AutoMigrate(&models.AttendanceQr{})
	if err != nil {
		log.Fatalf("Failed to migrate the table AttendanceQr: %v", err)
	}

	err = db.AutoMigrate(&models.StudentAttendance{})
	if err != nil {
		log.Fatalf("Failed to migrate the table StudentAttendance: %v", err)
	}

	err = db.AutoMigrate(&models.SectionSchedule{})
	if err != nil {
		log.Fatalf("Failed to migrate the table SectionSchedule: %v", err)
	}

	err = db.AutoMigrate(&models.StudentSectionEnrollment{})
	if err != nil {
		log.Fatalf("Failed to migrate the table StudentSectionEnrollment: %v", err)
	}

	log.Println("Database migration completed successfully")
}

// HealthCheck - To verify the database connectivity
func HealthCheck() error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}
	return sqlDB.Ping()
}

// CloseDB - Closes the database connection
func CloseDB() error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("%S: %w", utils.SQL_CONNECTION_ERROR, err)
	}
	return sqlDB.Close()
}

// GetDB - Getter function for getting the DB type variable
func GetDB() *gorm.DB {
	if db == nil {
		initDB()
		autoMigrate()
	}
	return db
}
