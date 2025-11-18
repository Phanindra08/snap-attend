package config

import (
	"fmt"
	"log"
	"os"

	"github.com/phanindra08/snap-attend/internal/shared/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port      int    `mapstructure:"port"`
		JwtSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"server"`

	Database struct {
		Host         string `mapstructure:"host"`
		Port         int    `mapstructure:"Port"`
		Username     string `mapstructure:"username"`
		Password     string `mapstructure:"password"`
		DatabaseName string `mapstructure:"database-name"`
		Schema       string `mapstructure:"schema"`
		SslMode      string `mapstructure:"ssl-mode"`
	} `mapstructure:"database"`
}

var appConfiguration *Config

// GetConnectionStringForDB returns the formatted PostGre SQL connection string
func (configuration *Config) GetConnectionStringForDB() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		configuration.Database.Host, configuration.Database.Port,
		configuration.Database.Username, configuration.Database.Password,
		configuration.Database.DatabaseName, configuration.Database.SslMode,
		configuration.Database.Schema)
}

func (configuration *Config) validateConfig() {
	required := []struct {
		value     string
		name      string
		minLength int
	}{
		{appConfiguration.Server.JwtSecret, utils.VALIDATING_JWT_SECRET, 32},
		{appConfiguration.Database.Host, utils.VALIDATING_HOST, 1},
		{appConfiguration.Database.Username, utils.VALIDATING_DB_USERNAME, 1},
		{appConfiguration.Database.Password, utils.VALIDATING_PASSWORD, 5},
		{appConfiguration.Database.DatabaseName, utils.VALIDATING_DB_NAME, 1},
		{appConfiguration.Database.Schema, utils.VALIDATING_SCHEMA, 1},
	}

	// Validating the length of the DB configurations and JWT Secret
	for _, req := range required {
		if len(req.value) < req.minLength {
			log.Fatalf("Configuration error: %s must be at least %d characters long.", req.name, req.minLength)
		}
	}

	// Validate if the DB Port is a valid number
	if configuration.Database.Port <= utils.MIN_TCP_PORT || configuration.Database.Port > utils.MAX_TCP_PORT {
		log.Fatalf("Configuration error: DB PORT must be between %d and %d", utils.MIN_TCP_PORT, utils.MAX_TCP_PORT)
	}

	// Validate if the server Port is a valid number
	if configuration.Server.Port <= utils.MIN_TCP_PORT || configuration.Server.Port > utils.MAX_TCP_PORT {
		log.Fatalf("Configuration error: Server PORT must be between %d and %d", utils.MIN_TCP_PORT, utils.MAX_TCP_PORT)
	}

	// Warn if the server Port is using a privileged Port
	if configuration.Server.Port < 1024 {
		log.Printf("Warning: Server Port %d is in the privileged range (1 - %d). Binding to this Port may require root privileges.", configuration.Server.Port, utils.MAX_PRIVILEGED_PORT)
	}
}

// GetServerPort - returns the Port the application is running
func (configuration *Config) GetServerPort() int {
	return configuration.Server.Port
}

// GetConfig - Getter function for getting the Config type variable
func GetConfig() *Config {
	if appConfiguration == nil {
		loadConfiguration()
	}
	return appConfiguration
}

// loadConfiguration - To load the Config variable
func loadConfiguration() {

	// Retrieving the environment variable SNAP_ATTEND_ENV for figuring out the config file to be used
	environment := os.Getenv(utils.ENVIRONMENT_VARIABLE)
	if environment == "" {
		log.Fatalf("Environment variable '%s' is not defined", utils.ENVIRONMENT_VARIABLE)
	}

	// Loading the information related to configuration file
	viper.SetConfigName(fmt.Sprintf("%s.%s", utils.CONFIG_NAME, environment))
	viper.SetConfigType(utils.CONFIG_TYPE)
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading the configuration file: %v", err)
	}

	appConfiguration = &Config{} // Initializing the struct
	if err := viper.Unmarshal(appConfiguration); err != nil {
		log.Fatalf("Unable to parse and decode into struct: %v", err)
	}

	appConfiguration.validateConfig()
}
