package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/notification"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler exposes notification HTTP endpoints.
type NotificationHandler struct {
	svc *notification.Service
}

// NewNotificationHandler constructs NotificationHandler.
func NewNotificationHandler(svc *notification.Service) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// Create godoc
// @Summary Create notification
// @Description Staff only. Optional metadata JSON may include "telegram_chat_id" for Telegram outbound when a bot is wired.
// @Tags notifications
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateNotificationRequest true "Notification"
// @Success 201 {object} response.Envelope{data=dto.NotificationResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/notifications [post]
func (h *NotificationHandler) Create(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateNotificationRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	nt, err := domain.ParseNotificationType(req.Type)
	if err != nil {
		response.Error(c, apperror.Validation("type", "Unknown notification type; use a supported type such as payment_reminder, grade_posted, or system"))
		return
	}
	out, err := h.svc.Create(c.Request.Context(), role, notification.CreateInput{
		UserID:   req.UserID,
		Type:     nt,
		Title:    req.Title,
		Body:     req.Body,
		Metadata: req.Metadata,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.NotificationResponseFrom(out))
}

// List godoc
// @Summary List notifications
// @Description Recipients see their own feed. Staff may pass user_id to list another user's notifications.
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param user_id query string false "Target user UUID (staff only)"
// @Param unread_only query string false "true or 1 for unread only"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Success 200 {object} response.Envelope{data=dto.NotificationListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	role, uid, err := RequireNotificationInboxActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.NotificationListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	var target *uuid.UUID
	if q.UserID != "" {
		parsed := uuid.MustParse(q.UserID)
		target = &parsed
	}
	unreadOnly := q.UnreadOnly == "true" || q.UnreadOnly == "1"
	page, pageSize := q.Page, q.PageSize
	items, total, err := h.svc.List(c.Request.Context(), role, uid, target, unreadOnly, page, pageSize)
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
	response.JSON(c, http.StatusOK, dto.NotificationListResponseFrom(items, total, pg, sz))
}

// GetByID godoc
// @Summary Get notification by ID
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.NotificationResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/notifications/{id} [get]
func (h *NotificationHandler) GetByID(c *gin.Context) {
	role, uid, err := RequireNotificationInboxActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.Get(c.Request.Context(), role, uid, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.NotificationResponseFrom(out))
}

// MarkRead godoc
// @Summary Mark notification as read
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.NotificationResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/notifications/{id}/read [patch]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	role, uid, err := RequireNotificationInboxActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.MarkRead(c.Request.Context(), role, uid, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.NotificationResponseFrom(out))
}

// Delete godoc
// @Summary Delete notification
// @Tags notifications
// @Security BearerAuth
// @Produce json
// @Param id path string true "Notification ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/notifications/{id} [delete]
func (h *NotificationHandler) Delete(c *gin.Context) {
	role, uid, err := RequireNotificationInboxActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), role, uid, id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Notification deleted successfully"})
}
