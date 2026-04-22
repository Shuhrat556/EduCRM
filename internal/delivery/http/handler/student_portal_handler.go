package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/usecase/studentportal"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StudentPortalHandler serves read-only student self-service APIs.
type StudentPortalHandler struct {
	svc *studentportal.Service
}

// NewStudentPortalHandler constructs StudentPortalHandler.
func NewStudentPortalHandler(svc *studentportal.Service) *StudentPortalHandler {
	return &StudentPortalHandler{svc: svc}
}

func (h *StudentPortalHandler) requireStudent(c *gin.Context) (uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return uuid.Nil, err
	}
	if role != domain.RoleStudent {
		return uuid.Nil, apperror.Forbidden("student portal only")
	}
	return uid, nil
}

// MyGrades GET /api/v1/student/grades
func (h *StudentPortalHandler) MyGrades(c *gin.Context) {
	uid, err := h.requireStudent(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	rows, err := h.svc.MyGrades(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GradeListResponseFrom(rows))
}

// MySchedule GET /api/v1/student/schedule
func (h *StudentPortalHandler) MySchedule(c *gin.Context) {
	uid, err := h.requireStudent(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	rows, err := h.svc.MySchedule(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.ScheduleListResponseFrom(rows))
}

// MyAttendance GET /api/v1/student/attendance
func (h *StudentPortalHandler) MyAttendance(c *gin.Context) {
	uid, err := h.requireStudent(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	rows, err := h.svc.MyAttendance(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AttendanceListResponseFrom(rows))
}
