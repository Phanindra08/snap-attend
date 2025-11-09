package utils

import "gorm.io/gorm/logger"

const (
	ENVIRONMENT_VARIABLE = "SNAP_ATTEND_ENV"
	CONFIG_NAME          = "config"
	CONFIG_TYPE          = "yaml"
)

const (
	VALIDATING_JWT_SECRET  = "JWT_Secret"
	VALIDATING_HOST        = "DB Host"
	VALIDATING_DB_USERNAME = "Username"
	VALIDATING_PASSWORD    = "DB Password"
	VALIDATING_DB_NAME     = "DB Name"
	VALIDATING_SCHEMA      = "Schema"
)

const (
	MIN_TCP_PORT              = 0
	MAX_TCP_PORT              = 65535
	MAX_PRIVILEGED_PORT       = 1024
	MAX_DB_RETRIES            = 5
	MAX_IDLE_CONNECTIONS      = 5
	MAX_OPEN_CONNECTIONS      = 100
	CONNECTIONS_MAX_IDLE_TIME = 5
)

const (
	DB_LOG_LEVEL = logger.Info
)

const (
	SQL_CONNECTION_ERROR = "failed to get database instance"
)
