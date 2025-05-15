package main

import (
	"fmt"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/logger"
	"backend/internal/models"
	"backend/internal/router"
)

func main() {
	appConfig := config.LoadConfig()

	logger.LogInfo("Starting BPV-Matcher backend...")

	db, err := database.Initialize()
	if err != nil {
		logger.LogFatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	logger.LogInfo("Running database migrations...")
	err = db.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Coordinator{},
		&models.Mentor{},
		&models.Company{},
	)
	if err != nil {
		logger.LogFatal(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	r := router.NewRouter()

	logger.LogInfof("Server running on port %s", appConfig.ServerPort)
	if err := r.Run(":" + appConfig.ServerPort); err != nil {
		logger.LogFatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
