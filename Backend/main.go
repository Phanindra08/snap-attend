package main

import (
	"github.com/phanindra08/snap-attend/internal/shared/config"
	"github.com/phanindra08/snap-attend/internal/shared/database"
)

func main() {
	config.GetConfig()
	database.GetDB()
	defer database.CloseDB()

	err := database.HealthCheck()
	if err != nil {
		return
	}
}
