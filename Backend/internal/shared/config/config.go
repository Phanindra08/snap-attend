package config

import (
	"fmt"
	"log"
	"os"

	"github.com/phanindra08/snap-attend/internal/shared/utils"
	"github.com/spf13/viper"
)

type Config struct {
	server struct {
		port      int    `mapstructure:"port"`
		jwtSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"server"`

	database struct {
		host         string `mapstructure:"host"`
		port         int    `mapstructure:"port"`
		username     string `mapstructure:"username"`
		password     string `mapstructure:"password"`
		databaseName string `mapstructure:"database-name"`
		schema       string `mapstructure:"schema"`
		sslMode      string `mapstructure:"ssl-mode"`
	} `mapstructure:"database"`
}

var appConfiguration *Config

// GetConnectionStringForDB returns the formatted PostGre SQL connection string
func (configuration *Config) GetConnectionStringForDB() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		configuration.database.host, configuration.database.port,
		configuration.database.username, configuration.database.password,
		configuration.database.databaseName, configuration.database.sslMode)
}

func (configuration *Config) validateConfig() {
	required := []struct {
		value     string
		name      string
		minLength int
	}{
		{appConfiguration.server.jwtSecret, utils.VALIDATING_JWT_SECRET, 32},
		{appConfiguration.database.host, utils.VALIDATING_HOST, 32},
		{appConfiguration.database.username, utils.VALIDATING_DB_USERNAME, 1},
		{appConfiguration.database.password, utils.VALIDATING_PASSWORD, 5},
		{appConfiguration.database.databaseName, utils.VALIDATING_DB_NAME, 1},
		{appConfiguration.database.schema, utils.VALIDATING_SCHEMA, 1},
	}

	// Validating the length of the DB configurations and JWT Secret
	for _, req := range required {
		if len(req.value) < req.minLength {
			log.Fatalf("Configuration error: %s must be at least %d characters long.", req.name, req.minLength)
		}
	}

	// Validate if the DB port is a valid number
	if configuration.database.port <= utils.MIN_TCP_PORT || configuration.database.port > utils.MAX_TCP_PORT {
		log.Fatalf("Configuration error: DB PORT must be between %d and %d", utils.MIN_TCP_PORT, utils.MAX_TCP_PORT)
	}

	// Validate if the server port is a valid number
	if configuration.server.port <= utils.MIN_TCP_PORT || configuration.server.port > utils.MAX_TCP_PORT {
		log.Fatalf("Configuration error: Server PORT must be between %d and %d", utils.MIN_TCP_PORT, utils.MAX_TCP_PORT)
	}

	// Warn if the server port is using a privileged port
	if configuration.server.port < 1024 {
		log.Printf("Warning: Server port %d is in the privileged range (1 - %d). Binding to this port may require root privileges.", configuration.server.port, utils.MAX_PRIVILEGED_PORT)
	}
}

// GetServerPort - returns the Port the application is running
func (configuration *Config) GetServerPort() int {
	return configuration.server.port
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

	if err := viper.Unmarshal(appConfiguration); err != nil {
		log.Fatalf("Unable to parse and decode into struct: %v", err)
	}

	appConfiguration.validateConfig()
}
