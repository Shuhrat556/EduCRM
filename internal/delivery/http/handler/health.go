package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/database"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler exposes liveness and readiness checks.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler constructs a health handler with injected DB (may be nil for tests without DB).
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Live godoc
// @Summary Liveness probe
// @Description Returns OK if the process is running. Does not check dependencies.
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope "data contains status"
// @Router /health [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"status": "ok"})
}

// Ready godoc
// @Summary Readiness probe
// @Description Returns OK if the application and database are ready to serve traffic.
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope "data contains status and database"
// @Failure 503 {object} response.Envelope "Service unavailable when database ping fails"
// @Router /api/v1/health [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.db == nil {
		response.JSON(c, http.StatusOK, gin.H{
			"status":   "ok",
			"database": "not_configured",
		})
		return
	}
	if err := database.Ping(h.db); err != nil {
		response.JSON(c, http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"database": "down",
		})
		return
	}
	response.JSON(c, http.StatusOK, gin.H{
		"status":   "ok",
		"database": "up",
	})
}
