package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/usecase/teacherportal"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TeacherPortalHandler serves teacher-scoped read APIs (assignments, roster, schedule).
type TeacherPortalHandler struct {
	svc *teacherportal.Service
}

// NewTeacherPortalHandler constructs TeacherPortalHandler.
func NewTeacherPortalHandler(svc *teacherportal.Service) *TeacherPortalHandler {
	return &TeacherPortalHandler{svc: svc}
}

func (h *TeacherPortalHandler) requireTeacher(c *gin.Context) (uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return uuid.Nil, err
	}
	if role != domain.RoleTeacher {
		return uuid.Nil, apperror.Forbidden("teacher portal only")
	}
	return uid, nil
}

// ListAssignments GET /api/v1/teacher/assignments
func (h *TeacherPortalHandler) ListAssignments(c *gin.Context) {
	uid, err := h.requireTeacher(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	rows, err := h.svc.ListAssignments(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"items": rows})
}

// ListStudentsQuery optional group_id
type ListStudentsQuery struct {
	GroupID *uuid.UUID `form:"group_id"`
}

// ListAssignedStudents GET /api/v1/teacher/students
func (h *TeacherPortalHandler) ListAssignedStudents(c *gin.Context) {
	uid, err := h.requireTeacher(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q ListStudentsQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	rows, err := h.svc.ListAssignedStudents(c.Request.Context(), uid, q.GroupID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, gin.H{"items": rows})
}

// MySchedule GET /api/v1/teacher/schedule
func (h *TeacherPortalHandler) MySchedule(c *gin.Context) {
	uid, err := h.requireTeacher(c)
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
