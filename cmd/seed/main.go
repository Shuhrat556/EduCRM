// Command seed creates the initial super_admin user when one does not already exist.
//
// Required env: SEED_SUPER_ADMIN_EMAIL, SEED_SUPER_ADMIN_PASSWORD (min 12 characters recommended).
// Optional: SEED_SUPER_ADMIN_PHONE
package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/internal/database"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres"
	authsvc "github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/google/uuid"
)

func main() {
	log.SetFlags(0)
	ctx := context.Background()

	email := strings.TrimSpace(os.Getenv("SEED_SUPER_ADMIN_EMAIL"))
	password := os.Getenv("SEED_SUPER_ADMIN_PASSWORD")
	phoneRaw := strings.TrimSpace(os.Getenv("SEED_SUPER_ADMIN_PHONE"))

	if email == "" || password == "" {
		log.Fatal("SEED_SUPER_ADMIN_EMAIL and SEED_SUPER_ADMIN_PASSWORD are required")
	}
	if !strings.Contains(email, "@") {
		log.Fatal("SEED_SUPER_ADMIN_EMAIL must be an email address")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
	}()

	repo := postgres.NewUserRepository(db)
	normEmail := domain.NormalizeEmail(&email)
	if normEmail == nil {
		log.Fatal("invalid email")
	}
	taken, err := repo.EmailTaken(ctx, strings.TrimSpace(strings.ToLower(*normEmail)), nil)
	if err != nil {
		log.Fatalf("email check: %v", err)
	}
	if taken {
		log.Println("seed: user with this email already exists; skipping")
		return
	}

	hash, err := authsvc.HashPassword(password)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}

	now := time.Now().UTC()
	u := &domain.User{
		ID:           uuid.New(),
		Email:        normEmail,
		Phone:        domain.NormalizePhone(nonEmptyPtr(phoneRaw)),
		PasswordHash: hash,
		Role:         domain.RoleSuperAdmin,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := repo.Create(ctx, u); err != nil {
		if err == repository.ErrDuplicate {
			log.Println("seed: duplicate user; skipping")
			return
		}
		log.Fatalf("create user: %v", err)
	}
	log.Printf("seed: created super_admin id=%s email=%s", u.ID, *u.Email)
}

func nonEmptyPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
