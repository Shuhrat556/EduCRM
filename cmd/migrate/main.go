// Command migrate applies versioned SQL migrations from ./migrations (or MIGRATIONS_PATH)
// against the database defined by DB_* environment variables.
//
// Usage: migrate <up|down|version|drop>
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/educrm/educrm-backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up|down|version|drop>")
	}
	cmd := strings.ToLower(strings.TrimSpace(os.Args[1]))

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	migrationsDir := strings.TrimSpace(os.Getenv("MIGRATIONS_PATH"))
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}
	abs, err := filepath.Abs(migrationsDir)
	if err != nil {
		log.Fatalf("migrations path: %v", err)
	}
	srcURL := fileURL(abs)
	dbURL := cfg.DB.PostgresURL()

	m, err := migrate.New(srcURL, dbURL)
	if err != nil {
		log.Fatalf("migrate init: %v", err)
	}
	defer func() { _, _ = m.Close() }()

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("up: %v", err)
		}
		log.Println("migrate: up OK")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("down: %v", err)
		}
		log.Println("migrate: down OK")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			if err == migrate.ErrNilVersion {
				log.Println("version: no migrations applied")
				return
			}
			log.Fatalf("version: %v", err)
		}
		log.Printf("version: %d dirty=%v", v, dirty)
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("drop: %v", err)
		}
		log.Println("migrate: drop OK")
	default:
		log.Fatalf("unknown command %q (use up, down, version, drop)", cmd)
	}
}

func fileURL(abs string) string {
	s := filepath.ToSlash(abs)
	if strings.HasPrefix(s, "/") {
		return "file://" + s
	}
	return "file:///" + s
}
