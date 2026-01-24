package handler

import (
	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/ncr"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/qc"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/traceability"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QCHandler handles QC-related requests
type QCHandler struct {
	getCheckpointsUC    *qc.GetCheckpointsUseCase
	createInspectionUC  *qc.CreateInspectionUseCase
	getInspectionUC     *qc.GetInspectionUseCase
	listInspectionsUC   *qc.ListInspectionsUseCase
	approveInspectionUC *qc.ApproveInspectionUseCase
}

// NewQCHandler creates a new QCHandler
func NewQCHandler(
	getCheckpointsUC *qc.GetCheckpointsUseCase,
	createInspectionUC *qc.CreateInspectionUseCase,
	getInspectionUC *qc.GetInspectionUseCase,
	listInspectionsUC *qc.ListInspectionsUseCase,
	approveInspectionUC *qc.ApproveInspectionUseCase,
) *QCHandler {
	return &QCHandler{
		getCheckpointsUC:    getCheckpointsUC,
		createInspectionUC:  createInspectionUC,
		getInspectionUC:     getInspectionUC,
		listInspectionsUC:   listInspectionsUC,
		approveInspectionUC: approveInspectionUC,
	}
}

// GetCheckpoints gets all QC checkpoints
func (h *QCHandler) GetCheckpoints(c *gin.Context) {
	checkpoints, err := h.getCheckpointsUC.Execute(c.Request.Context())
	if err != nil {
		internalError(c, err.Error())
		return
	}
	success(c, checkpoints)
}

// CreateInspection creates a new QC inspection
func (h *QCHandler) CreateInspection(c *gin.Context) {
	var req dto.CreateInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	var items []qc.CreateInspectionItemInput
	for i, item := range req.Items {
		items = append(items, qc.CreateInspectionItemInput{
			ItemNumber:    i + 1,
			TestName:      item.TestName,
			TestMethod:    item.TestMethod,
			Specification: item.Specification,
			TargetValue:   item.TargetValue,
			MinValue:      item.MinValue,
			MaxValue:      item.MaxValue,
			ActualValue:   item.ActualValue,
			UOM:           item.UOM,
			Result:        entity.ItemResult(item.Result),
			Notes:         item.Notes,
		})
	}

	input := qc.CreateInspectionInput{
		InspectionType:    entity.CheckpointType(req.InspectionType),
		CheckpointID:      req.CheckpointID,
		ReferenceType:     entity.ReferenceType(req.ReferenceType),
		ReferenceID:       req.ReferenceID,
		ProductID:         req.ProductID,
		MaterialID:        req.MaterialID,
		LotID:             req.LotID,
		LotNumber:         req.LotNumber,
		InspectedQuantity: req.InspectedQuantity,
		SampleSize:        req.SampleSize,
		InspectorID:       userID,
		InspectorName:     req.InspectorName,
		Items:             items,
	}

	result, err := h.createInspectionUC.Execute(c.Request.Context(), input)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	created(c, result)
}

// GetInspection gets a QC inspection by ID
func (h *QCHandler) GetInspection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid inspection ID")
		return
	}

	result, err := h.getInspectionUC.Execute(c.Request.Context(), id)
	if err != nil {
		notFound(c, "Inspection not found")
		return
	}

	success(c, result)
}

// ListInspections lists QC inspections
func (h *QCHandler) ListInspections(c *gin.Context) {
	filter := repository.QCFilter{
		Page:     getPageFromQuery(c),
		PageSize: getPageSizeFromQuery(c),
	}

	if t := c.Query("type"); t != "" {
		ct := entity.CheckpointType(t)
		filter.InspectionType = &ct
	}
	if r := c.Query("result"); r != "" {
		ir := entity.InspectionResult(r)
		filter.Result = &ir
	}

	inspections, total, err := h.listInspectionsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	successWithMeta(c, inspections, newMeta(filter.Page, filter.PageSize, total))
}

// ApproveInspection approves or rejects an inspection
func (h *QCHandler) ApproveInspection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid inspection ID")
		return
	}

	var req dto.ApproveInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	input := qc.ApproveInspectionInput{
		InspectionID:     id,
		Result:           entity.InspectionResult(req.Result),
		AcceptedQuantity: req.AcceptedQuantity,
		RejectedQuantity: req.RejectedQuantity,
		ApproverID:       userID,
		Notes:            req.Notes,
	}

	result, err := h.approveInspectionUC.Execute(c.Request.Context(), input)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, result)
}

// NCRHandler handles NCR-related requests
type NCRHandler struct {
	createNCRUC *ncr.CreateNCRUseCase
	getNCRUC    *ncr.GetNCRUseCase
	listNCRsUC  *ncr.ListNCRsUseCase
	closeNCRUC  *ncr.CloseNCRUseCase
}

// NewNCRHandler creates a new NCRHandler
func NewNCRHandler(
	createNCRUC *ncr.CreateNCRUseCase,
	getNCRUC *ncr.GetNCRUseCase,
	listNCRsUC *ncr.ListNCRsUseCase,
	closeNCRUC *ncr.CloseNCRUseCase,
) *NCRHandler {
	return &NCRHandler{
		createNCRUC: createNCRUC,
		getNCRUC:    getNCRUC,
		listNCRsUC:  listNCRsUC,
		closeNCRUC:  closeNCRUC,
	}
}

// CreateNCR creates a new NCR
func (h *NCRHandler) CreateNCR(c *gin.Context) {
	var req dto.CreateNCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	severity := entity.NCRSeverityMedium
	if req.Severity != "" {
		severity = entity.NCRSeverity(req.Severity)
	}

	input := ncr.CreateNCRInput{
		NCType:           entity.NCType(req.NCType),
		Severity:         severity,
		ReferenceType:    req.ReferenceType,
		ReferenceID:      req.ReferenceID,
		ProductID:        req.ProductID,
		MaterialID:       req.MaterialID,
		LotID:            req.LotID,
		LotNumber:        req.LotNumber,
		Description:      req.Description,
		QuantityAffected: req.QuantityAffected,
		UOMID:            req.UOMID,
		ImmediateAction:  req.ImmediateAction,
		CreatedBy:        userID,
	}

	result, err := h.createNCRUC.Execute(c.Request.Context(), input)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	created(c, result)
}

// GetNCR gets an NCR by ID
func (h *NCRHandler) GetNCR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid NCR ID")
		return
	}

	result, err := h.getNCRUC.Execute(c.Request.Context(), id)
	if err != nil {
		notFound(c, "NCR not found")
		return
	}

	success(c, result)
}

// ListNCRs lists NCRs
func (h *NCRHandler) ListNCRs(c *gin.Context) {
	filter := repository.NCRFilter{
		Page:     getPageFromQuery(c),
		PageSize: getPageSizeFromQuery(c),
	}

	if status := c.Query("status"); status != "" {
		s := entity.NCRStatus(status)
		filter.Status = &s
	}
	if severity := c.Query("severity"); severity != "" {
		sev := entity.NCRSeverity(severity)
		filter.Severity = &sev
	}

	ncrs, total, err := h.listNCRsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	successWithMeta(c, ncrs, newMeta(filter.Page, filter.PageSize, total))
}

// CloseNCR closes an NCR
func (h *NCRHandler) CloseNCR(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid NCR ID")
		return
	}

	var req dto.CloseNCRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	var disposition *entity.Disposition
	if req.Disposition != "" {
		d := entity.Disposition(req.Disposition)
		disposition = &d
	}

	input := ncr.CloseNCRInput{
		NCRID:            id,
		RootCause:        req.RootCause,
		CorrectiveAction: req.CorrectiveAction,
		PreventiveAction: req.PreventiveAction,
		Disposition:      disposition,
		DispositionQty:   req.DispositionQty,
		ClosureNotes:     req.ClosureNotes,
		ClosedBy:         userID,
	}

	result, err := h.closeNCRUC.Execute(c.Request.Context(), input)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, result)
}

// TraceHandler handles traceability requests
type TraceHandler struct {
	traceBackwardUC *traceability.TraceBackwardUseCase
	traceForwardUC  *traceability.TraceForwardUseCase
}

// NewTraceHandler creates a new TraceHandler
func NewTraceHandler(
	traceBackwardUC *traceability.TraceBackwardUseCase,
	traceForwardUC *traceability.TraceForwardUseCase,
) *TraceHandler {
	return &TraceHandler{
		traceBackwardUC: traceBackwardUC,
		traceForwardUC:  traceForwardUC,
	}
}

// TraceBackward traces backward from product lot to material lots
func (h *TraceHandler) TraceBackward(c *gin.Context) {
	id, err := uuid.Parse(c.Param("lot_id"))
	if err != nil {
		badRequest(c, "Invalid lot ID")
		return
	}

	result, err := h.traceBackwardUC.Execute(c.Request.Context(), id)
	if err != nil {
		notFound(c, "Traceability data not found")
		return
	}

	success(c, result)
}

// TraceForward traces forward from material lot to product lots
func (h *TraceHandler) TraceForward(c *gin.Context) {
	id, err := uuid.Parse(c.Param("lot_id"))
	if err != nil {
		badRequest(c, "Invalid lot ID")
		return
	}

	result, err := h.traceForwardUC.Execute(c.Request.Context(), id)
	if err != nil {
		notFound(c, "Traceability data not found")
		return
	}

	success(c, result)
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns service health
func (h *HealthHandler) Health(c *gin.Context) {
	success(c, gin.H{
		"status":  "healthy",
		"service": "manufacturing-service",
	})
}

// Ready returns readiness status
func (h *HealthHandler) Ready(c *gin.Context) {
	success(c, gin.H{"status": "ready"})
}

// Live returns liveness status
func (h *HealthHandler) Live(c *gin.Context) {
	success(c, gin.H{"status": "live"})
}
