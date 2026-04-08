package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Env      string `env:"APP_ENV" envDefault:"development"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// ShutdownTimeout is the maximum time to wait for in-flight HTTP requests on SIGINT/SIGTERM.
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`

	// AutoMigrate runs GORM AutoMigrate on API startup. Disable in production when using SQL migrations only.
	AutoMigrate bool `env:"AUTO_MIGRATE" envDefault:"true"`

	// SwaggerEnabled exposes /swagger when true. If ENABLE_SWAGGER is unset, defaults to on except in production.
	SwaggerEnabled bool `env:"ENABLE_SWAGGER"`

	HTTP      HTTPConfig
	DB        DatabaseConfig
	JWT       JWTConfig
	Storage   StorageConfig
	AI        AIConfig `envPrefix:"AI_"`
	CORS      CORSConfig
	RateLimit RateLimitConfig
	LogHTTP   LogHTTPConfig
}

// HTTPConfig defines the HTTP server settings.
type HTTPConfig struct {
	Host         string        `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port         string        `env:"HTTP_PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"15s"`
	IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	// TrustedProxies is a comma-separated list of CIDRs or IPs trusted for X-Forwarded-For (empty = gin default).
	TrustedProxies string `env:"HTTP_TRUSTED_PROXIES"`
}

// DatabaseConfig defines PostgreSQL connection settings.
type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"localhost"`
	Port            string        `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER" envDefault:"educrm"`
	Password        string        `env:"DB_PASSWORD" envDefault:"educrm"`
	Name            string        `env:"DB_NAME" envDefault:"educrm"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"5m"`
	// DebugSQL enables verbose GORM SQL logs (development only; avoid in production).
	DebugSQL bool `env:"DB_DEBUG_SQL" envDefault:"false"`
}

// JWTConfig defines signing settings for access JWTs (refresh tokens are opaque, stored in DB).
type JWTConfig struct {
	Secret            string        `env:"JWT_SECRET" envDefault:"change-me-in-production"`
	AccessExpiration  time.Duration `env:"JWT_ACCESS_EXPIRATION" envDefault:"15m"`
	RefreshExpiration time.Duration `env:"JWT_REFRESH_EXPIRATION" envDefault:"168h"`
	Issuer            string        `env:"JWT_ISSUER" envDefault:"educrm"`
}

// CORSConfig controls browser cross-origin access.
type CORSConfig struct {
	// AllowedOrigins is a comma-separated list (e.g. https://app.example.com,https://admin.example.com). Empty in development defaults to *.
	AllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`
	// AllowCredentials sets Access-Control-Allow-Credentials (requires non-wildcard origins).
	AllowCredentials bool `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
}

// RateLimitConfig is an in-process rate limit foundation (per client IP). Use a shared store in multi-instance deployments.
type RateLimitConfig struct {
	Enabled bool    `env:"RATE_LIMIT_ENABLED" envDefault:"false"`
	RPS     float64 `env:"RATE_LIMIT_RPS" envDefault:"100"`
	Burst   int     `env:"RATE_LIMIT_BURST" envDefault:"200"`
}

// LogHTTPConfig tunes request access logs.
type LogHTTPConfig struct {
	// SkipPaths is a comma-separated list of path prefixes excluded from request logs (e.g. /health,/api/v1/health).
	SkipPaths string `env:"LOG_HTTP_SKIP_PATHS" envDefault:"/health,/api/v1/health"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env config: %w", err)
	}
	applyDerivedDefaults(cfg)
	return cfg, nil
}

func applyDerivedDefaults(cfg *Config) {
	if _, set := os.LookupEnv("ENABLE_SWAGGER"); !set {
		cfg.SwaggerEnabled = !cfg.IsProduction()
	}
}

// DSN returns the PostgreSQL connection string for GORM/pg.
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// Addr returns the HTTP listen address.
func (c HTTPConfig) Addr() string {
	return c.Host + ":" + c.Port
}

// TrustedProxyList splits HTTP_TRUSTED_PROXIES into a slice for gin.Engine.SetTrustedProxies.
func (c HTTPConfig) TrustedProxyList() []string {
	s := strings.TrimSpace(c.TrustedProxies)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// LogHTTPSkipPrefixes returns path prefixes to omit from access logs.
func (c *Config) LogHTTPSkipPrefixes() []string {
	if c == nil {
		return nil
	}
	s := strings.TrimSpace(c.LogHTTP.SkipPaths)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
