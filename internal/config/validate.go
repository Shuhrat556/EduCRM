package config

import (
	"fmt"
	"strings"
)

var weakJWTSecrets = map[string]struct{}{
	"change-me-in-production":   {},
	"dev-only-change-me":        {},
	"your-256-bit-secret":       {},
	"secret":                    {},
}

// IsProduction reports whether APP_ENV is production.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(strings.TrimSpace(c.Env), "production")
}

// IsStaging reports whether APP_ENV is staging.
func (c *Config) IsStaging() bool {
	return strings.EqualFold(strings.TrimSpace(c.Env), "staging")
}

// IsDevelopment reports development-like environments (default when unset).
func (c *Config) IsDevelopment() bool {
	if c.IsProduction() || c.IsStaging() {
		return false
	}
	return true
}

// ValidateForAPI enforces settings required for serving HTTP in production/staging.
// CLI tools (migrate, seed) should not call this; they only need database config.
func (c *Config) ValidateForAPI() error {
	if c.IsProduction() {
		return c.validateProductionAPI()
	}
	if c.IsStaging() {
		return c.validateStagingAPI()
	}
	return nil
}

func (c *Config) validateProductionAPI() error {
	sec := strings.TrimSpace(c.JWT.Secret)
	if len(sec) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters in production")
	}
	if isWeakJWTSecret(sec) {
		return fmt.Errorf("JWT_SECRET must not use a default or placeholder value in production")
	}
	if strings.EqualFold(strings.TrimSpace(c.DB.SSLMode), "disable") {
		return fmt.Errorf("DB_SSLMODE must not be disable in production (use require, verify-ca, or verify-full)")
	}
	origins := c.parseCORSOrigins()
	if len(origins) == 0 {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS must be set in production (comma-separated origins, no bare wildcard)")
	}
	for _, o := range origins {
		if o == "*" {
			return fmt.Errorf("CORS_ALLOWED_ORIGINS cannot be * in production")
		}
	}
	if c.CORS.AllowCredentials {
		for _, o := range origins {
			if o == "*" {
				return fmt.Errorf("CORS_ALLOW_CREDENTIALS cannot be used with wildcard origins")
			}
		}
	}
	return nil
}

func (c *Config) validateStagingAPI() error {
	sec := strings.TrimSpace(c.JWT.Secret)
	if len(sec) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters in staging")
	}
	if isWeakJWTSecret(sec) {
		return fmt.Errorf("JWT_SECRET must not use a default or placeholder value in staging")
	}
	return nil
}

func isWeakJWTSecret(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	_, ok := weakJWTSecrets[strings.ToLower(s)]
	return ok
}
