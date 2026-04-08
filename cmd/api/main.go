package main

import (
	"log"
	"os"

	"github.com/educrm/educrm-backend/internal/app"
	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/pkg/logger"
)

// @title EduCRM API
// @version 1.0
// @description REST API for EduCRM: JWT auth (access + refresh), users, teachers, rooms, groups, schedules, attendance, grades, payments, files, notifications, dashboard, and AI analytics. Success responses use a JSON envelope: `{"success":true,"data":...}`. Errors use `{"success":false,"error":{"code","message","kind"}}`.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@educrm.local

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	if err := cfg.ValidateForAPI(); err != nil {
		log.Fatalf("config validation: %v", err)
	}

	logr := logger.New(cfg.LogLevel, cfg.Env)

	application, err := app.New(cfg, logr)
	if err != nil {
		logr.Error("app_init", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := application.Close(); err != nil {
			logr.Error("app_close", "error", err)
		}
	}()

	if err := application.Run(); err != nil {
		logr.Error("app_run", "error", err)
		os.Exit(1)
	}
}
