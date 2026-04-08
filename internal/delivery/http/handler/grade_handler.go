package handler

import (
	"math"
	"net/http"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/grade"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GradeHandler exposes grade HTTP endpoints.
type GradeHandler struct {
	svc *grade.Service
}

// NewGradeHandler constructs GradeHandler.
func NewGradeHandler(svc *grade.Service) *GradeHandler {
	return &GradeHandler{svc: svc}
}

// Create godoc
// @Summary Create weekly grade
// @Description One row per student, group, teacher, calendar week (Monday UTC), and grade_type. teacher_evaluation: admin or assigned teacher. student_evaluation: admin, the student (self), or assigned teacher.
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateGradeRequest true "Grade"
// @Success 201 {object} response.Envelope{data=dto.GradeResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/grades [post]
func (h *GradeHandler) Create(c *gin.Context) {
	role, actorID, err := RequireGradesActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateGradeRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if math.IsNaN(req.GradeValue) || math.IsInf(req.GradeValue, 0) {
		response.Error(c, apperror.Validation("grade_value", "grade_value must be a finite number"))
		return
	}
	gt, err := domain.ParseGradeType(req.GradeType)
	if err != nil {
		response.Error(c, apperror.Validation("grade_type", "grade_type must be teacher_evaluation or student_evaluation"))
		return
	}
	in := grade.CreateInput{
		StudentID:  req.StudentID,
		GroupID:    req.GroupID,
		GradeType:  gt,
		GradeValue: req.GradeValue,
		Comment:    req.Comment,
	}
	if req.WeekOf != nil && *req.WeekOf != "" {
		d, err := dto.ParseISODate(*req.WeekOf)
		if err != nil {
			response.Error(c, apperror.Validation("week_of", "Use date format YYYY-MM-DD"))
			return
		}
		in.WeekOf = &d
	}
	if req.GradedAt != nil && *req.GradedAt != "" {
		t, err := time.Parse(time.RFC3339, *req.GradedAt)
		if err != nil {
			response.Error(c, apperror.Validation("graded_at", "Use RFC3339 datetime (e.g. 2025-04-08T15:04:05Z)"))
			return
		}
		in.GradedAt = &t
	}
	out, err := h.svc.Create(c.Request.Context(), role, actorID, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.GradeResponseFrom(out))
}

// List godoc
// @Summary List grades
// @Description Exactly one of student_id or group_id. Students may only use student_id for themselves. Teachers need user_teacher_links.
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param student_id query string false "Student user UUID"
// @Param group_id query string false "Group UUID"
// @Success 200 {object} response.Envelope{data=dto.GradeListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/grades [get]
func (h *GradeHandler) List(c *gin.Context) {
	role, actorID, err := RequireGradesActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.GradeListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	var f grade.ListFilter
	if q.StudentID != "" {
		id := uuid.MustParse(q.StudentID)
		f.StudentID = &id
	}
	if q.GroupID != "" {
		id := uuid.MustParse(q.GroupID)
		f.GroupID = &id
	}
	rows, err := h.svc.List(c.Request.Context(), role, actorID, f)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GradeListResponseFrom(rows))
}

// GetByID godoc
// @Summary Get grade by ID
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Grade ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.GradeResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/grades/{id} [get]
func (h *GradeHandler) GetByID(c *gin.Context) {
	role, actorID, err := RequireGradesActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), role, actorID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GradeResponseFrom(out))
}

// Update godoc
// @Summary Update grade
// @Tags grades
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Grade ID (UUID)"
// @Param body body dto.UpdateGradeRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.GradeResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/grades/{id} [patch]
func (h *GradeHandler) Update(c *gin.Context) {
	role, actorID, err := RequireGradesActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateGradeRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := grade.UpdateInput{GradeValue: req.GradeValue, Comment: req.Comment}
	if req.GradedAt != nil && *req.GradedAt != "" {
		t, err := time.Parse(time.RFC3339, *req.GradedAt)
		if err != nil {
			response.Error(c, apperror.Validation("graded_at", "Use RFC3339 datetime"))
			return
		}
		in.GradedAt = &t
	}
	if in.GradeValue == nil && in.Comment == nil && in.GradedAt == nil {
		response.Error(c, apperror.Validation("body", "Provide at least one of grade_value, comment, or graded_at"))
		return
	}
	if in.GradeValue != nil && (math.IsNaN(*in.GradeValue) || math.IsInf(*in.GradeValue, 0)) {
		response.Error(c, apperror.Validation("grade_value", "grade_value must be a finite number"))
		return
	}
	out, err := h.svc.Update(c.Request.Context(), role, actorID, id, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GradeResponseFrom(out))
}

// Delete godoc
// @Summary Delete grade
// @Tags grades
// @Security BearerAuth
// @Produce json
// @Param id path string true "Grade ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/grades/{id} [delete]
func (h *GradeHandler) Delete(c *gin.Context) {
	role, actorID, err := RequireGradesActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), role, actorID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Grade deleted successfully"})
}
