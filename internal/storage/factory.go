package storage

import (
	"fmt"
	"strings"

	"github.com/educrm/educrm-backend/internal/config"
)

// NewFromConfig returns the configured Provider. S3/MinIO are reserved for a
// future implementation; until then use provider "local" or register metadata
// via HTTP after uploading to object storage yourself.
func NewFromConfig(cfg config.StorageConfig) (Provider, error) {
	p := strings.ToLower(strings.TrimSpace(cfg.Provider))
	switch p {
	case "local", "":
		return NewLocal(cfg.LocalDir, cfg.PublicBaseURL)
	case "s3", "minio":
		return nil, fmt.Errorf("storage provider %q is not implemented yet; use STORAGE_PROVIDER=local or POST /files/register after external upload", cfg.Provider)
	default:
		return nil, fmt.Errorf("unknown STORAGE_PROVIDER %q", cfg.Provider)
	}
}
