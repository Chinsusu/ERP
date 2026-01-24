package dto

import (
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
	"github.com/google/uuid"
)

// CreatePRRequest represents request to create a PR
type CreatePRRequest struct {
	RequiredDate  string              `json:"required_date" binding:"required"`
	Priority      string              `json:"priority"`
	Justification string              `json:"justification"`
	Notes         string              `json:"notes"`
	Items         []PRLineItemRequest `json:"items" binding:"required,min=1"`
}

// PRLineItemRequest represents a line item in the request
type PRLineItemRequest struct {
	MaterialID     string  `json:"material_id" binding:"required"`
	MaterialCode   string  `json:"material_code"`
	MaterialName   string  `json:"material_name"`
	Quantity       float64 `json:"quantity" binding:"required,gt=0"`
	UOMCode        string  `json:"uom_code"`
	UnitPrice      float64 `json:"unit_price"`
	Specifications string  `json:"specifications"`
}

// ApprovePRRequest represents request to approve/reject a PR
type ApprovePRRequest struct {
	Notes  string `json:"notes"`
	Reason string `json:"reason"`
}

// PRResponse represents the PR response
type PRResponse struct {
	ID            string              `json:"id"`
	PRNumber      string              `json:"pr_number"`
	PRDate        string              `json:"pr_date"`
	RequiredDate  string              `json:"required_date"`
	Priority      string              `json:"priority"`
	Status        string              `json:"status"`
	RequesterID   string              `json:"requester_id"`
	DepartmentID  *string             `json:"department_id,omitempty"`
	Justification string              `json:"justification,omitempty"`
	Notes         string              `json:"notes,omitempty"`
	Currency      string              `json:"currency"`
	TotalAmount   float64             `json:"total_amount"`
	ApprovalLevel int                 `json:"approval_level"`
	ApprovedBy    *string             `json:"approved_by,omitempty"`
	ApprovedAt    *string             `json:"approved_at,omitempty"`
	RejectedBy    *string             `json:"rejected_by,omitempty"`
	RejectedAt    *string             `json:"rejected_at,omitempty"`
	RejectReason  string              `json:"reject_reason,omitempty"`
	ConvertedToPO *string             `json:"converted_to_po,omitempty"`
	LineItems     []PRLineItemResponse `json:"line_items,omitempty"`
	Approvals     []PRApprovalResponse `json:"approvals,omitempty"`
	CreatedAt     string              `json:"created_at"`
	UpdatedAt     string              `json:"updated_at"`
}

// PRLineItemResponse represents a line item response
type PRLineItemResponse struct {
	ID             string  `json:"id"`
	LineNumber     int     `json:"line_number"`
	MaterialID     string  `json:"material_id"`
	MaterialCode   string  `json:"material_code"`
	MaterialName   string  `json:"material_name"`
	Quantity       float64 `json:"quantity"`
	UOMCode        string  `json:"uom_code"`
	UnitPrice      float64 `json:"unit_price"`
	Currency       string  `json:"currency"`
	LineTotal      float64 `json:"line_total"`
	Specifications string  `json:"specifications,omitempty"`
}

// PRApprovalResponse represents an approval record response
type PRApprovalResponse struct {
	ID            string `json:"id"`
	ApproverID    string `json:"approver_id"`
	ApprovalLevel int    `json:"approval_level"`
	Action        string `json:"action"`
	Notes         string `json:"notes,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// ToPRResponse converts entity to response
func ToPRResponse(pr *entity.PurchaseRequisition) *PRResponse {
	resp := &PRResponse{
		ID:            pr.ID.String(),
		PRNumber:      pr.PRNumber,
		PRDate:        pr.PRDate.Format("2006-01-02"),
		RequiredDate:  pr.RequiredDate.Format("2006-01-02"),
		Priority:      string(pr.Priority),
		Status:        string(pr.Status),
		RequesterID:   pr.RequesterID.String(),
		Justification: pr.Justification,
		Notes:         pr.Notes,
		Currency:      pr.Currency,
		TotalAmount:   pr.TotalAmount,
		ApprovalLevel: int(getApprovalLevelInt(pr.ApprovalLevel)),
		RejectReason:  pr.RejectionReason,
		CreatedAt:     pr.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     pr.UpdatedAt.Format(time.RFC3339),
	}

	if pr.DepartmentID != nil {
		deptID := pr.DepartmentID.String()
		resp.DepartmentID = &deptID
	}
	if pr.ApprovedBy != nil {
		approvedBy := pr.ApprovedBy.String()
		resp.ApprovedBy = &approvedBy
	}
	if pr.ApprovedAt != nil {
		approvedAt := pr.ApprovedAt.Format(time.RFC3339)
		resp.ApprovedAt = &approvedAt
	}
	if pr.RejectedBy != nil {
		rejectedBy := pr.RejectedBy.String()
		resp.RejectedBy = &rejectedBy
	}
	if pr.RejectedAt != nil {
		rejectedAt := pr.RejectedAt.Format(time.RFC3339)
		resp.RejectedAt = &rejectedAt
	}
	if pr.POID != nil {
		poID := pr.POID.String()
		resp.ConvertedToPO = &poID
	}

	for _, item := range pr.LineItems {
		resp.LineItems = append(resp.LineItems, PRLineItemResponse{
			ID:             item.ID.String(),
			LineNumber:     item.LineNumber,
			MaterialID:     item.MaterialID.String(),
			MaterialCode:   item.MaterialCode,
			MaterialName:   item.MaterialName,
			Quantity:       item.Quantity,
			UOMCode:        item.UOMCode,
			UnitPrice:      item.UnitPrice,
			Currency:       item.Currency,
			LineTotal:      item.LineTotal,
			Specifications: item.Specifications,
		})
	}

	for _, approval := range pr.Approvals {
		resp.Approvals = append(resp.Approvals, PRApprovalResponse{
			ID:            approval.ID.String(),
			ApproverID:    approval.ApproverID.String(),
			ApprovalLevel: int(getApprovalLevelInt(approval.ApprovalLevel)),
			Action:        approval.Action,
			Notes:         approval.Notes,
			CreatedAt:     approval.CreatedAt.Format(time.RFC3339),
		})
	}

	return resp
}

// getApprovalLevelInt converts ApprovalLevel to int
func getApprovalLevelInt(level entity.ApprovalLevel) int {
	switch level {
	case entity.ApprovalLevelAuto:
		return 0
	case entity.ApprovalLevelDeptManager:
		return 1
	case entity.ApprovalLevelProcurementManager:
		return 2
	case entity.ApprovalLevelCFO:
		return 3
	default:
		return 0
	}
}

// ToPRListResponse converts list of entities to response
func ToPRListResponse(prs []*entity.PurchaseRequisition) []*PRResponse {
	var responses []*PRResponse
	for _, pr := range prs {
		responses = append(responses, ToPRResponse(pr))
	}
	return responses
}

// ParseMaterialID parses material ID string to UUID
func ParseMaterialID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
