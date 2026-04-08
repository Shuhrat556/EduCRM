package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/usecase/payment"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentHandler exposes payment HTTP endpoints.
type PaymentHandler struct {
	svc *payment.Service
}

// NewPaymentHandler constructs PaymentHandler.
func NewPaymentHandler(svc *payment.Service) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

// Create godoc
// @Summary Create payment
// @Description Amounts are minor units (e.g. cents). Only super_admin may set is_free or discount_amount_minor &gt; 0. Multiple rows per student/group/month_for are allowed for partial payments.
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreatePaymentRequest true "Payment"
// @Success 201 {object} response.Envelope{data=dto.PaymentResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreatePaymentRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	st, err := domain.ParsePaymentStatus(req.Status)
	if err != nil {
		response.Error(c, apperror.Validation("status", "Status must be paid_full, paid_partial, unpaid, or overdue"))
		return
	}
	pt, err := domain.ParsePaymentType(req.PaymentType)
	if err != nil {
		response.Error(c, apperror.Validation("payment_type", "Invalid payment_type value"))
		return
	}
	mf, err := dto.ParseMonthFor(req.MonthFor)
	if err != nil {
		response.Error(c, apperror.Validation("month_for", "Use YYYY-MM or YYYY-MM-DD"))
		return
	}
	in := payment.CreateInput{
		StudentID:           req.StudentID,
		GroupID:             req.GroupID,
		AmountMinor:         req.AmountMinor,
		Status:              st,
		MonthFor:            mf,
		PaymentType:         pt,
		Comment:             req.Comment,
		IsFree:              req.IsFree,
		DiscountAmountMinor: req.DiscountAmountMinor,
	}
	if req.PaymentDate != nil && *req.PaymentDate != "" {
		d, err := dto.ParseISODate(*req.PaymentDate)
		if err != nil {
			response.Error(c, apperror.Validation("payment_date", "Use date format YYYY-MM-DD"))
			return
		}
		in.PaymentDate = &d
	}
	out, err := h.svc.Create(c.Request.Context(), role, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.PaymentResponseFrom(out))
}

// List godoc
// @Summary List payments (staff)
// @Description Paginated staff view with optional filters.
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Param q query string false "Search in comment"
// @Param student_id query string false "Filter by student UUID"
// @Param group_id query string false "Filter by group UUID"
// @Param month_for query string false "Filter billed month (YYYY-MM or YYYY-MM-DD)"
// @Param status query string false "paid_full | paid_partial | unpaid | overdue"
// @Param payment_type query string false "monthly_tuition | partial_payment | adjustment | other"
// @Param is_free query string false "true | false | 1 | 0"
// @Success 200 {object} response.Envelope{data=dto.PaymentListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments [get]
func (h *PaymentHandler) List(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.PaymentListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	params := repository.PaymentListParams{
		Search:   q.Q,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if q.StudentID != "" {
		sid := uuid.MustParse(q.StudentID)
		params.StudentID = &sid
	}
	if q.GroupID != "" {
		gid := uuid.MustParse(q.GroupID)
		params.GroupID = &gid
	}
	if q.MonthFor != "" {
		mf, err := dto.ParseMonthFor(q.MonthFor)
		if err != nil {
			response.Error(c, apperror.Validation("month_for", "Use YYYY-MM or YYYY-MM-DD"))
			return
		}
		t := mf
		params.MonthFor = &t
	}
	if q.Status != "" {
		st, err := domain.ParsePaymentStatus(q.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Invalid status filter"))
			return
		}
		params.Status = &st
	}
	if q.PaymentType != "" {
		pt, err := domain.ParsePaymentType(q.PaymentType)
		if err != nil {
			response.Error(c, apperror.Validation("payment_type", "Invalid payment_type filter"))
			return
		}
		params.PaymentType = &pt
	}
	switch q.IsFree {
	case "true", "1":
		v := true
		params.IsFree = &v
	case "false", "0":
		v := false
		params.IsFree = &v
	}
	out, err := h.svc.List(c.Request.Context(), role, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.PaymentListResponseFrom(out))
}

// History godoc
// @Summary Student payment history
// @Description Students see their own history. Staff must pass student_id.
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Param student_id query string false "Required for admin/super_admin"
// @Param page query int false "Page" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Success 200 {object} response.Envelope{data=dto.PaymentListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments/history [get]
func (h *PaymentHandler) History(c *gin.Context) {
	role, uid, err := RequirePaymentsReadActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.PaymentHistoryQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	var sid *uuid.UUID
	if q.StudentID != "" {
		parsed := uuid.MustParse(q.StudentID)
		sid = &parsed
	}
	out, err := h.svc.History(c.Request.Context(), role, uid, sid, q.Page, q.PageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.PaymentListResponseFrom(out))
}

// GetByID godoc
// @Summary Get payment by ID
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Payment ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.PaymentResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetByID(c *gin.Context) {
	role, uid, err := RequirePaymentsReadActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), role, uid, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.PaymentResponseFrom(out))
}

// Update godoc
// @Summary Update payment
// @Description Only super_admin may set is_free true or discount_amount_minor &gt; 0.
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Payment ID (UUID)"
// @Param body body dto.UpdatePaymentRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.PaymentResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments/{id} [patch]
func (h *PaymentHandler) Update(c *gin.Context) {
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
	var req dto.UpdatePaymentRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := payment.UpdateInput{}
	if req.AmountMinor != nil {
		in.AmountMinor = req.AmountMinor
	}
	if req.Status != nil {
		st, err := domain.ParsePaymentStatus(*req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be paid_full, paid_partial, unpaid, or overdue"))
			return
		}
		in.Status = &st
	}
	if req.PaymentDate != nil && *req.PaymentDate != "" {
		d, err := dto.ParseISODate(*req.PaymentDate)
		if err != nil {
			response.Error(c, apperror.Validation("payment_date", "Use date format YYYY-MM-DD"))
			return
		}
		in.PaymentDate = &d
	}
	if req.PaymentType != nil {
		pt, err := domain.ParsePaymentType(*req.PaymentType)
		if err != nil {
			response.Error(c, apperror.Validation("payment_type", "Invalid payment_type value"))
			return
		}
		in.PaymentType = &pt
	}
	if req.Comment != nil {
		in.Comment = req.Comment
	}
	if req.IsFree != nil {
		in.IsFree = req.IsFree
	}
	if req.DiscountAmountMinor != nil {
		in.DiscountAmountMinor = req.DiscountAmountMinor
	}
	out, err := h.svc.Update(c.Request.Context(), role, id, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.PaymentResponseFrom(out))
}

// Delete godoc
// @Summary Soft-delete payment
// @Description Row is soft-deleted; history remains in the database.
// @Tags payments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Payment ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/payments/{id} [delete]
func (h *PaymentHandler) Delete(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Payment deleted successfully"})
}
