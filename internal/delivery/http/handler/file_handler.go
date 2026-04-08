package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/file"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileHandler handles file uploads and metadata.
type FileHandler struct {
	svc     *file.Service
	maxMem  int64
	maxBody int64
}

// NewFileHandler constructs FileHandler (maxMultipartMemory bounds buffering for ParseMultipartForm).
func NewFileHandler(svc *file.Service, st config.StorageConfig) *FileHandler {
	max := st.MaxUploadBytes
	if max <= 0 {
		max = 10 << 20
	}
	mem := max + (1 << 20)
	if mem < 8<<20 {
		mem = 8 << 20
	}
	return &FileHandler{svc: svc, maxMem: mem, maxBody: max}
}

// Upload godoc
// @Summary Upload file (multipart)
// @Description Form fields: file (binary), owner_type (student_photo | teacher_photo | document), owner_id (UUID). Stored via configured storage provider; metadata row is created.
// @Tags files
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param file formData file true "File"
// @Param owner_type formData string true "student_photo | teacher_photo | document"
// @Param owner_id formData string true "Owner UUID"
// @Success 201 {object} response.Envelope{data=dto.FileMetadataResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 413 {object} response.Envelope "Payload too large"
// @Failure 500 {object} response.Envelope
// @Router /api/v1/files [post]
func (h *FileHandler) Upload(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := c.Request.ParseMultipartForm(h.maxMem); err != nil {
		response.Error(c, apperror.Validation("multipart", "Could not parse multipart form; check Content-Type and field names"))
		return
	}
	ownerTypeStr := c.PostForm("owner_type")
	ownerIDStr := c.PostForm("owner_id")
	ot, err := domain.ParseFileOwnerType(ownerTypeStr)
	if err != nil {
		response.Error(c, apperror.Validation("owner_type", "owner_type must be student_photo, teacher_photo, or document"))
		return
	}
	ownerID, err := uuid.Parse(ownerIDStr)
	if err != nil {
		response.Error(c, apperror.Validation("owner_id", "owner_id must be a valid UUID"))
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.Validation("file", "A file part named \"file\" is required"))
		return
	}
	if fh.Size > h.maxBody {
		response.Error(c, apperror.Validation("file", "File exceeds the maximum allowed upload size"))
		return
	}
	src, err := fh.Open()
	if err != nil {
		response.Error(c, apperror.Internal("open upload").Wrap(err))
		return
	}
	defer src.Close()

	mt := contentType(fh)
	out, err := h.svc.Upload(c.Request.Context(), role, file.UploadInput{
		OwnerType: ot,
		OwnerID:   ownerID,
		FileName:  fh.Filename,
		MimeType:  mt,
		Size:      fh.Size,
		Body:      src,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.FileMetadataResponseFrom(out))
}

func contentType(fh *multipart.FileHeader) string {
	if fh.Header != nil {
		if v := fh.Header.Get("Content-Type"); v != "" {
			return v
		}
	}
	return "application/octet-stream"
}

// Register godoc
// @Summary Register file metadata (external storage)
// @Description Use after uploading to S3/MinIO (or when STORAGE_PROVIDER is not local). Optional storage_key enables Delete to remove the remote object later.
// @Tags files
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.RegisterFileMetadataRequest true "Metadata"
// @Success 201 {object} response.Envelope{data=dto.FileMetadataResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/files/register [post]
func (h *FileHandler) Register(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.RegisterFileMetadataRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	ot, err := domain.ParseFileOwnerType(req.OwnerType)
	if err != nil {
		response.Error(c, apperror.Validation("owner_type", "owner_type must be student_photo, teacher_photo, or document"))
		return
	}
	out, err := h.svc.RegisterMetadata(c.Request.Context(), role, file.RegisterInput{
		OwnerType:  ot,
		OwnerID:    req.OwnerID,
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		MimeType:   req.MimeType,
		Size:       req.Size,
		StorageKey: req.StorageKey,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.FileMetadataResponseFrom(out))
}

// List godoc
// @Summary List file metadata by owner
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param owner_type query string true "student_photo | teacher_photo | document"
// @Param owner_id query string true "Owner UUID"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Success 200 {object} response.Envelope{data=dto.FileMetadataListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/files [get]
func (h *FileHandler) List(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.FileListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	ot, err := domain.ParseFileOwnerType(q.OwnerType)
	if err != nil {
		response.Error(c, apperror.Validation("owner_type", "owner_type must be student_photo, teacher_photo, or document"))
		return
	}
	ownerID := uuid.MustParse(q.OwnerID)
	page, pageSize := q.Page, q.PageSize
	items, total, err := h.svc.List(c.Request.Context(), role, ot, ownerID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	pg := page
	if pg < 1 {
		pg = 1
	}
	sz := pageSize
	if sz < 1 {
		sz = 20
	}
	if sz > 100 {
		sz = 100
	}
	response.JSON(c, http.StatusOK, dto.FileMetadataListResponseFrom(items, total, pg, sz))
}

// GetByID godoc
// @Summary Get file metadata by ID
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param id path string true "File metadata ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.FileMetadataResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/files/{id} [get]
func (h *FileHandler) GetByID(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.Get(c.Request.Context(), role, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.FileMetadataResponseFrom(out))
}

// Delete godoc
// @Summary Delete file metadata and stored object
// @Description Removes DB row; deletes blob via storage provider unless the row was register-only (no object key).
// @Tags files
// @Security BearerAuth
// @Produce json
// @Param id path string true "File metadata ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/files/{id} [delete]
func (h *FileHandler) Delete(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), role, id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "File deleted successfully"})
}
