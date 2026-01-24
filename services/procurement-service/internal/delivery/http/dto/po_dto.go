package dto

import (
	"time"

	"github.com/erp-cosmetics/procurement-service/internal/domain/entity"
)

// ConvertPRToPORequest represents request to convert PR to PO
type ConvertPRToPORequest struct {
	SupplierID           string `json:"supplier_id" binding:"required"`
	SupplierCode         string `json:"supplier_code"`
	SupplierName         string `json:"supplier_name"`
	DeliveryAddress      string `json:"delivery_address"`
	DeliveryTerms        string `json:"delivery_terms"`
	PaymentTerms         string `json:"payment_terms"`
	ExpectedDeliveryDate string `json:"expected_delivery_date"`
	Notes                string `json:"notes"`
}

// ConfirmPORequest represents request to confirm a PO
type ConfirmPORequest struct {
	Notes string `json:"notes"`
}

// CancelPORequest represents request to cancel a PO
type CancelPORequest struct {
	Reason string `json:"reason" binding:"required"`
}

// AmendPORequest represents request to amend a PO
type AmendPORequest struct {
	FieldChanged string `json:"field_changed" binding:"required"`
	OldValue     string `json:"old_value"`
	NewValue     string `json:"new_value"`
	Reason       string `json:"reason"`
}

// POResponse represents the PO response
type POResponse struct {
	ID                   string               `json:"id"`
	PONumber             string               `json:"po_number"`
	PODate               string               `json:"po_date"`
	PRID                 *string              `json:"pr_id,omitempty"`
	SupplierID           string               `json:"supplier_id"`
	SupplierCode         string               `json:"supplier_code"`
	SupplierName         string               `json:"supplier_name"`
	Status               string               `json:"status"`
	DeliveryAddress      string               `json:"delivery_address,omitempty"`
	DeliveryTerms        string               `json:"delivery_terms"`
	PaymentTerms         string               `json:"payment_terms"`
	ExpectedDeliveryDate *string              `json:"expected_delivery_date,omitempty"`
	Currency             string               `json:"currency"`
	SubTotal             float64              `json:"sub_total"`
	TaxAmount            float64              `json:"tax_amount"`
	DiscountAmount       float64              `json:"discount_amount"`
	GrandTotal           float64              `json:"grand_total"`
	AmendmentCount       int                  `json:"amendment_count"`
	Notes                string               `json:"notes,omitempty"`
	ConfirmedBy          *string              `json:"confirmed_by,omitempty"`
	ConfirmedAt          *string              `json:"confirmed_at,omitempty"`
	CancelledBy          *string              `json:"cancelled_by,omitempty"`
	CancelledAt          *string              `json:"cancelled_at,omitempty"`
	CancelReason         string               `json:"cancel_reason,omitempty"`
	LineItems            []POLineItemResponse `json:"line_items,omitempty"`
	Amendments           []POAmendmentResponse `json:"amendments,omitempty"`
	CreatedAt            string               `json:"created_at"`
	UpdatedAt            string               `json:"updated_at"`
}

// POLineItemResponse represents a PO line item response
type POLineItemResponse struct {
	ID             string  `json:"id"`
	LineNumber     int     `json:"line_number"`
	MaterialID     string  `json:"material_id"`
	MaterialCode   string  `json:"material_code"`
	MaterialName   string  `json:"material_name"`
	Quantity       float64 `json:"quantity"`
	ReceivedQty    float64 `json:"received_qty"`
	PendingQty     float64 `json:"pending_qty"`
	UOMCode        string  `json:"uom_code"`
	UnitPrice      float64 `json:"unit_price"`
	Currency       string  `json:"currency"`
	LineTotal      float64 `json:"line_total"`
	Status         string  `json:"status"`
	Specifications string  `json:"specifications,omitempty"`
}

// POAmendmentResponse represents a PO amendment response
type POAmendmentResponse struct {
	ID              string `json:"id"`
	AmendmentNumber int    `json:"amendment_number"`
	FieldChanged    string `json:"field_changed"`
	OldValue        string `json:"old_value"`
	NewValue        string `json:"new_value"`
	Reason          string `json:"reason,omitempty"`
	AmendedBy       string `json:"amended_by"`
	CreatedAt       string `json:"created_at"`
}

// POReceiptResponse represents a PO receipt response
type POReceiptResponse struct {
	ID              string  `json:"id"`
	POLineItemID    string  `json:"po_line_item_id"`
	GRNID           *string `json:"grn_id,omitempty"`
	GRNNumber       string  `json:"grn_number,omitempty"`
	ReceivedQty     float64 `json:"received_qty"`
	AcceptedQty     float64 `json:"accepted_qty"`
	RejectedQty     float64 `json:"rejected_qty"`
	ReceivedDate    string  `json:"received_date"`
	QCStatus        string  `json:"qc_status"`
	QCNotes         string  `json:"qc_notes,omitempty"`
}

// ToPOResponse converts entity to response
func ToPOResponse(po *entity.PurchaseOrder) *POResponse {
	resp := &POResponse{
		ID:             po.ID.String(),
		PONumber:       po.PONumber,
		PODate:         po.PODate.Format("2006-01-02"),
		SupplierID:     po.SupplierID.String(),
		SupplierCode:   po.SupplierCode,
		SupplierName:   po.SupplierName,
		Status:         string(po.Status),
		DeliveryAddress: po.DeliveryAddress,
		DeliveryTerms:  po.DeliveryTerms,
		PaymentTerms:   po.PaymentTerms,
		Currency:       po.Currency,
		SubTotal:       po.TotalAmount,
		TaxAmount:      po.TaxAmount,
		DiscountAmount: 0,
		GrandTotal:     po.GrandTotal,
		AmendmentCount: po.AmendmentCount,
		Notes:          po.Notes,
		CancelReason:   po.CancellationReason,
		CreatedAt:      po.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      po.UpdatedAt.Format(time.RFC3339),
	}

	if po.PRID != nil {
		prID := po.PRID.String()
		resp.PRID = &prID
	}
	if po.ExpectedDeliveryDate != nil {
		expectedDate := po.ExpectedDeliveryDate.Format("2006-01-02")
		resp.ExpectedDeliveryDate = &expectedDate
	}
	if po.ConfirmedBy != nil {
		confirmedBy := po.ConfirmedBy.String()
		resp.ConfirmedBy = &confirmedBy
	}
	if po.ConfirmedAt != nil {
		confirmedAt := po.ConfirmedAt.Format(time.RFC3339)
		resp.ConfirmedAt = &confirmedAt
	}
	if po.CancelledBy != nil {
		cancelledBy := po.CancelledBy.String()
		resp.CancelledBy = &cancelledBy
	}
	if po.CancelledAt != nil {
		cancelledAt := po.CancelledAt.Format(time.RFC3339)
		resp.CancelledAt = &cancelledAt
	}

	for _, item := range po.LineItems {
		resp.LineItems = append(resp.LineItems, POLineItemResponse{
			ID:             item.ID.String(),
			LineNumber:     item.LineNumber,
			MaterialID:     item.MaterialID.String(),
			MaterialCode:   item.MaterialCode,
			MaterialName:   item.MaterialName,
			Quantity:       item.Quantity,
			ReceivedQty:    item.ReceivedQty,
			PendingQty:     item.PendingQty,
			UOMCode:        item.UOMCode,
			UnitPrice:      item.UnitPrice,
			Currency:       item.Currency,
			LineTotal:      item.LineTotal,
			Status:         string(item.Status),
			Specifications: item.Specifications,
		})
	}

	for _, amend := range po.Amendments {
		resp.Amendments = append(resp.Amendments, POAmendmentResponse{
			ID:              amend.ID.String(),
			AmendmentNumber: amend.AmendmentNumber,
			FieldChanged:    amend.FieldChanged,
			OldValue:        amend.OldValue,
			NewValue:        amend.NewValue,
			Reason:          amend.Reason,
			AmendedBy:       amend.AmendedBy.String(),
			CreatedAt:       amend.CreatedAt.Format(time.RFC3339),
		})
	}

	return resp
}

// ToPOListResponse converts list of entities to response
func ToPOListResponse(pos []*entity.PurchaseOrder) []*POResponse {
	var responses []*POResponse
	for _, po := range pos {
		responses = append(responses, ToPOResponse(po))
	}
	return responses
}

// ToPOReceiptResponses converts receipts to responses
func ToPOReceiptResponses(receipts []*entity.POReceipt) []POReceiptResponse {
	var responses []POReceiptResponse
	for _, r := range receipts {
		resp := POReceiptResponse{
			ID:           r.ID.String(),
			POLineItemID: r.POLineItemID.String(),
			ReceivedQty:  r.ReceivedQty,
			AcceptedQty:  0,
			RejectedQty:  0,
			ReceivedDate: r.ReceivedDate.Format("2006-01-02"),
			QCStatus:     string(r.QCStatus),
			QCNotes:      r.QCNotes,
			GRNNumber:    r.GRNNumber,
		}
		if r.GRNID != nil {
			grnID := r.GRNID.String()
			resp.GRNID = &grnID
		}
		responses = append(responses, resp)
	}
	return responses
}
