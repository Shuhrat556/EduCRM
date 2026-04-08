package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/storage"
	"github.com/google/uuid"
)

// registerOnlyPrefix marks rows with no blob in the current Provider (metadata-only).
const registerOnlyPrefix = "register-only/"

// Service manages file metadata and writes through storage.Provider.
type Service struct {
	meta     repository.FileMetadataRepository
	blob     storage.Provider
	users    repository.UserRepository
	teachers repository.TeacherRepository
	maxBytes int64
}

// NewService constructs the file upload service.
func NewService(
	meta repository.FileMetadataRepository,
	blob storage.Provider,
	users repository.UserRepository,
	teachers repository.TeacherRepository,
	maxUploadBytes int64,
) *Service {
	if maxUploadBytes <= 0 {
		maxUploadBytes = 10 << 20
	}
	return &Service{
		meta:     meta,
		blob:     blob,
		users:    users,
		teachers: teachers,
		maxBytes: maxUploadBytes,
	}
}

// UploadInput is a decoded multipart upload.
type UploadInput struct {
	OwnerType domain.FileOwnerType
	OwnerID   uuid.UUID
	FileName  string
	MimeType  string
	Size      int64
	Body      io.Reader
}

// RegisterInput records metadata after an external upload (S3/MinIO, etc.).
type RegisterInput struct {
	OwnerType  domain.FileOwnerType
	OwnerID    uuid.UUID
	FileName   string
	FileURL    string
	MimeType   string
	Size       int64
	StorageKey *string // object key for future Delete via Provider; omit for metadata-only rows
}

// Upload validates the owner, stores bytes via Provider, then persists metadata.
func (s *Service) Upload(ctx context.Context, actorRole domain.Role, in UploadInput) (*domain.FileMetadata, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	if err := s.assertOwner(ctx, in.OwnerType, in.OwnerID); err != nil {
		return nil, err
	}
	if err := validateMime(in.OwnerType, in.MimeType); err != nil {
		return nil, err
	}
	data, err := io.ReadAll(io.LimitReader(in.Body, s.maxBytes+1))
	if err != nil {
		return nil, apperror.Internal("read upload").Wrap(err)
	}
	if int64(len(data)) > s.maxBytes {
		return nil, apperror.Validation("file", "file exceeds maximum allowed size")
	}
	if in.Size > 0 && int64(len(data)) != in.Size {
		return nil, apperror.Validation("file", "size does not match uploaded content")
	}
	id := uuid.New()
	key := buildStorageKey(in.OwnerType, in.OwnerID, id, in.FileName, in.MimeType)
	if err := s.blob.Put(ctx, key, bytes.NewReader(data), int64(len(data)), in.MimeType); err != nil {
		return nil, apperror.Internal("store file").Wrap(err)
	}
	url := s.blob.PublicURL(key)
	now := time.Now().UTC()
	row := &domain.FileMetadata{
		ID:         id,
		OwnerType:  in.OwnerType,
		OwnerID:    in.OwnerID,
		FileName:   sanitizeFileName(in.FileName),
		StorageKey: key,
		FileURL:    url,
		MimeType:   in.MimeType,
		Size:       int64(len(data)),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.meta.Create(ctx, row); err != nil {
		_ = s.blob.Delete(ctx, key)
		return nil, apperror.Internal("save metadata").Wrap(err)
	}
	return row, nil
}

// RegisterMetadata persists metadata without calling Provider.Put (e.g. client uploaded to S3).
func (s *Service) RegisterMetadata(ctx context.Context, actorRole domain.Role, in RegisterInput) (*domain.FileMetadata, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	if err := s.assertOwner(ctx, in.OwnerType, in.OwnerID); err != nil {
		return nil, err
	}
	if err := validateMime(in.OwnerType, in.MimeType); err != nil {
		return nil, err
	}
	if strings.TrimSpace(in.FileURL) == "" {
		return nil, apperror.Validation("file_url", "required")
	}
	if in.Size < 0 {
		return nil, apperror.Validation("size", "must be >= 0")
	}
	sk := ""
	if in.StorageKey != nil {
		sk = strings.TrimSpace(*in.StorageKey)
	}
	if sk == "" {
		sk = registerOnlyPrefix + uuid.New().String()
	} else if err := validateExternalStorageKey(sk); err != nil {
		return nil, apperror.Validation("storage_key", err.Error())
	}
	id := uuid.New()
	now := time.Now().UTC()
	row := &domain.FileMetadata{
		ID:         id,
		OwnerType:  in.OwnerType,
		OwnerID:    in.OwnerID,
		FileName:   sanitizeFileName(in.FileName),
		StorageKey: sk,
		FileURL:    strings.TrimSpace(in.FileURL),
		MimeType:   in.MimeType,
		Size:       in.Size,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.meta.Create(ctx, row); err != nil {
		return nil, apperror.Internal("save metadata").Wrap(err)
	}
	return row, nil
}

// Get returns one row visible to staff.
func (s *Service) Get(ctx context.Context, actorRole domain.Role, id uuid.UUID) (*domain.FileMetadata, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	row, err := s.meta.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load file metadata").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("file")
	}
	return row, nil
}

// List returns paginated metadata for an owner.
func (s *Service) List(ctx context.Context, actorRole domain.Role, ownerType domain.FileOwnerType, ownerID uuid.UUID, page, pageSize int) ([]domain.FileMetadata, int64, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, 0, err
	}
	if err := s.assertOwner(ctx, ownerType, ownerID); err != nil {
		return nil, 0, err
	}
	items, total, err := s.meta.List(ctx, repository.FileMetadataListParams{
		OwnerType: string(ownerType),
		OwnerID:   ownerID,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, 0, apperror.Internal("list files").Wrap(err)
	}
	return items, total, nil
}

// Delete removes metadata and the blob when storage key is managed by this app.
func (s *Service) Delete(ctx context.Context, actorRole domain.Role, id uuid.UUID) error {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return err
	}
	row, err := s.meta.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load file metadata").Wrap(err)
	}
	if row == nil {
		return apperror.NotFound("file")
	}
	if !strings.HasPrefix(row.StorageKey, registerOnlyPrefix) {
		if err := s.blob.Delete(ctx, row.StorageKey); err != nil {
			return apperror.Internal("delete file from storage").Wrap(err)
		}
	}
	if err := s.meta.Delete(ctx, id); err != nil {
		return apperror.Internal("delete metadata").Wrap(err)
	}
	return nil
}

func (s *Service) assertOwner(ctx context.Context, ot domain.FileOwnerType, ownerID uuid.UUID) error {
	switch ot {
	case domain.FileOwnerStudentPhoto:
		u, err := s.users.FindByID(ctx, ownerID)
		if err != nil {
			return apperror.Internal("load user").Wrap(err)
		}
		if u == nil || u.Role != domain.RoleStudent {
			return apperror.Validation("owner_id", "must be an existing student user")
		}
	case domain.FileOwnerTeacherPhoto:
		t, _, err := s.teachers.FindByID(ctx, ownerID)
		if err != nil {
			return apperror.Internal("load teacher").Wrap(err)
		}
		if t == nil {
			return apperror.Validation("owner_id", "must be an existing teacher")
		}
	case domain.FileOwnerDocument:
		// owner_id is an application-defined anchor (e.g. group or enrollment UUID).
	}
	return nil
}

func validateMime(ot domain.FileOwnerType, mimeType string) error {
	mt := strings.TrimSpace(strings.ToLower(mimeType))
	if mt == "" {
		return apperror.Validation("mime_type", "required")
	}
	switch ot {
	case domain.FileOwnerStudentPhoto, domain.FileOwnerTeacherPhoto:
		if !strings.HasPrefix(mt, "image/") {
			return apperror.Validation("mime_type", "must be an image type for photo owner_type")
		}
	case domain.FileOwnerDocument:
		if strings.HasPrefix(mt, "image/") || mt == "application/pdf" ||
			mt == "application/msword" ||
			mt == "application/vnd.openxmlformats-officedocument.wordprocessingml.document" ||
			mt == "application/vnd.ms-excel" ||
			mt == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" ||
			strings.HasPrefix(mt, "text/") {
			return nil
		}
		return apperror.Validation("mime_type", "unsupported document mime type")
	}
	return nil
}

func buildStorageKey(ot domain.FileOwnerType, ownerID, fileID uuid.UUID, fileName, mimeType string) string {
	ext := strings.ToLower(filepath.Ext(sanitizeFileName(fileName)))
	if ext == "" || ext == "." {
		ext = extFromMime(mimeType)
	}
	return string(ot) + "/" + ownerID.String() + "/" + fileID.String() + ext
}

func extFromMime(mimeType string) string {
	exts, _ := mime.ExtensionsByType(strings.TrimSpace(mimeType))
	if len(exts) > 0 {
		return exts[0]
	}
	return ".bin"
}

func sanitizeFileName(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "." || base == "/" || base == "" {
		return "upload"
	}
	return base
}

func validateExternalStorageKey(sk string) error {
	if strings.Contains(sk, "..") {
		return fmt.Errorf("invalid storage key")
	}
	return nil
}
