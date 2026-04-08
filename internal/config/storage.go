package config

// StorageConfig selects a blob backend and local filesystem layout.
// Swap STORAGE_PROVIDER to s3 or minio once an object-store implementation is wired.
type StorageConfig struct {
	Provider       string `env:"STORAGE_PROVIDER" envDefault:"local"` // local | s3 | minio (s3/minio: upload via /files/register until driver exists)
	LocalDir       string `env:"STORAGE_LOCAL_DIR" envDefault:"./data/uploads"`
	PublicBaseURL  string `env:"STORAGE_PUBLIC_BASE_URL" envDefault:"http://localhost:8080/static/files"`
	MaxUploadBytes int64  `env:"STORAGE_MAX_UPLOAD_BYTES" envDefault:"10485760"` // 10 MiB
}
