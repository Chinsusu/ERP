package handler

import (
	"strconv"

	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/erp-cosmetics/supplier-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/certification"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/evaluation"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CertificationHandler handles certification HTTP requests
type CertificationHandler struct {
	addCertUC        *certification.AddCertificationUseCase
	getCertsUC       *certification.GetCertificationsUseCase
	getExpiringUC    *certification.GetExpiringCertificationsUseCase
}

// NewCertificationHandler creates a new CertificationHandler
func NewCertificationHandler(
	addCertUC *certification.AddCertificationUseCase,
	getCertsUC *certification.GetCertificationsUseCase,
	getExpiringUC *certification.GetExpiringCertificationsUseCase,
) *CertificationHandler {
	return &CertificationHandler{
		addCertUC:     addCertUC,
		getCertsUC:    getCertsUC,
		getExpiringUC: getExpiringUC,
	}
}

// ListCertifications handles GET /api/v1/suppliers/:id/certifications
func (h *CertificationHandler) ListCertifications(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	certs, err := h.getCertsUC.Execute(c.Request.Context(), supplierID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.CertificationResponse
	for _, cert := range certs {
		items = append(items, *dto.ToCertificationResponse(cert))
	}

	response.Success(c, items)
}

// AddCertification handles POST /api/v1/suppliers/:id/certifications
func (h *CertificationHandler) AddCertification(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.CreateCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	ucReq := &certification.AddCertificationRequest{
		SupplierID:  supplierID,
		Type:        req.Type,
		CertNumber:  req.CertNumber,
		IssuingBody: req.IssuingBody,
		IssueDate:   req.IssueDate,
		ExpiryDate:  req.ExpiryDate,
		DocumentURL: req.DocumentURL,
		Notes:       req.Notes,
	}

	result, err := h.addCertUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ToCertificationResponse(result))
}

// GetExpiringCertifications handles GET /api/v1/certifications/expiring
func (h *CertificationHandler) GetExpiringCertifications(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	certs, err := h.getExpiringUC.Execute(c.Request.Context(), days)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.CertificationResponse
	for _, cert := range certs {
		items = append(items, *dto.ToCertificationResponse(cert))
	}

	response.Success(c, gin.H{
		"data":  items,
		"total": len(items),
		"days":  days,
	})
}

// EvaluationHandler handles evaluation HTTP requests
type EvaluationHandler struct {
	createEvalUC *evaluation.CreateEvaluationUseCase
	getEvalsUC   *evaluation.GetEvaluationsUseCase
}

// NewEvaluationHandler creates a new EvaluationHandler
func NewEvaluationHandler(
	createEvalUC *evaluation.CreateEvaluationUseCase,
	getEvalsUC *evaluation.GetEvaluationsUseCase,
) *EvaluationHandler {
	return &EvaluationHandler{
		createEvalUC: createEvalUC,
		getEvalsUC:   getEvalsUC,
	}
}

// ListEvaluations handles GET /api/v1/suppliers/:id/evaluations
func (h *EvaluationHandler) ListEvaluations(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	evals, err := h.getEvalsUC.Execute(c.Request.Context(), supplierID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.EvaluationResponse
	for _, e := range evals {
		items = append(items, *dto.ToEvaluationResponse(e))
	}

	response.Success(c, items)
}

// CreateEvaluation handles POST /api/v1/suppliers/:id/evaluations
func (h *EvaluationHandler) CreateEvaluation(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = uuid.New() // Fallback for testing
	}

	ucReq := &evaluation.CreateEvaluationRequest{
		SupplierID:            supplierID,
		EvaluationDate:        req.EvaluationDate,
		EvaluationPeriod:      req.EvaluationPeriod,
		QualityScore:          req.QualityScore,
		DeliveryScore:         req.DeliveryScore,
		PriceScore:            req.PriceScore,
		ServiceScore:          req.ServiceScore,
		DocumentationScore:    req.DocumentationScore,
		OnTimeDeliveryRate:    req.OnTimeDeliveryRate,
		QualityAcceptanceRate: req.QualityAcceptanceRate,
		LeadTimeAdherence:     req.LeadTimeAdherence,
		Strengths:             req.Strengths,
		Weaknesses:            req.Weaknesses,
		ActionItems:           req.ActionItems,
		EvaluatedBy:           userID,
	}

	result, err := h.createEvalUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ToEvaluationResponse(result))
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"service": "supplier-service",
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	response.Success(c, gin.H{"status": "ready"})
}
