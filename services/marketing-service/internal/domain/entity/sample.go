package entity

import (
	"time"

	"github.com/google/uuid"
)

// SampleRequestStatus represents sample request workflow status
type SampleRequestStatus string

const (
	SampleStatusDraft            SampleRequestStatus = "DRAFT"
	SampleStatusPendingApproval  SampleRequestStatus = "PENDING_APPROVAL"
	SampleStatusApproved         SampleRequestStatus = "APPROVED"
	SampleStatusRejected         SampleRequestStatus = "REJECTED"
	SampleStatusShipped          SampleRequestStatus = "SHIPPED"
	SampleStatusDelivered        SampleRequestStatus = "DELIVERED"
	SampleStatusFeedbackReceived SampleRequestStatus = "FEEDBACK_RECEIVED"
)

// SampleRequest represents a request to send product samples to KOL
type SampleRequest struct {
	ID              uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RequestNumber   string              `json:"request_number" gorm:"uniqueIndex;size:50"`
	KOLID           uuid.UUID           `json:"kol_id" gorm:"type:uuid;not null"`
	KOL             *KOL                `json:"kol,omitempty" gorm:"foreignKey:KOLID"`
	CampaignID      *uuid.UUID          `json:"campaign_id" gorm:"type:uuid"`
	Campaign        *Campaign           `json:"campaign,omitempty" gorm:"foreignKey:CampaignID"`
	CollaborationID *uuid.UUID          `json:"collaboration_id" gorm:"type:uuid"`

	RequestDate   time.Time `json:"request_date"`
	RequestReason string    `json:"request_reason"`

	// Delivery
	DeliveryAddress string `json:"delivery_address"`
	RecipientName   string `json:"recipient_name" gorm:"size:200"`
	RecipientPhone  string `json:"recipient_phone" gorm:"size:50"`

	// Expectations
	ExpectedPostDate *time.Time `json:"expected_post_date"`
	ExpectedReach    int        `json:"expected_reach"`

	// Value
	TotalItems int     `json:"total_items"`
	TotalValue float64 `json:"total_value" gorm:"type:decimal(18,2)"`

	Status SampleRequestStatus `json:"status" gorm:"size:50;default:DRAFT"`

	// Approval
	ApprovedBy      *uuid.UUID `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt      *time.Time `json:"approved_at"`
	RejectionReason string     `json:"rejection_reason"`

	Notes string `json:"notes"`

	Items []SampleItem `json:"items,omitempty" gorm:"foreignKey:SampleRequestID"`

	CreatedBy *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (SampleRequest) TableName() string {
	return "sample_requests"
}

// CanBeApproved checks if request can be approved
func (s *SampleRequest) CanBeApproved() bool {
	return s.Status == SampleStatusDraft || s.Status == SampleStatusPendingApproval
}

// Approve approves the sample request
func (s *SampleRequest) Approve(approverID uuid.UUID) {
	now := time.Now()
	s.Status = SampleStatusApproved
	s.ApprovedBy = &approverID
	s.ApprovedAt = &now
	s.UpdatedAt = now
}

// Reject rejects the sample request
func (s *SampleRequest) Reject(reason string) {
	s.Status = SampleStatusRejected
	s.RejectionReason = reason
	s.UpdatedAt = time.Now()
}

// MarkShipped marks request as shipped
func (s *SampleRequest) MarkShipped() {
	s.Status = SampleStatusShipped
	s.UpdatedAt = time.Now()
}

// MarkDelivered marks request as delivered
func (s *SampleRequest) MarkDelivered() {
	s.Status = SampleStatusDelivered
	s.UpdatedAt = time.Now()
}

// MarkFeedbackReceived marks feedback as received
func (s *SampleRequest) MarkFeedbackReceived() {
	s.Status = SampleStatusFeedbackReceived
	s.UpdatedAt = time.Now()
}

// CalculateTotals recalculates totals from line items
func (s *SampleRequest) CalculateTotals() {
	var totalItems int
	var totalValue float64
	for _, item := range s.Items {
		totalItems += item.Quantity
		totalValue += item.TotalValue
	}
	s.TotalItems = totalItems
	s.TotalValue = totalValue
}

// SampleItem represents a product in a sample request
type SampleItem struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SampleRequestID uuid.UUID  `json:"sample_request_id" gorm:"type:uuid;not null"`
	LineNumber      int        `json:"line_number"`
	
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	ProductCode string    `json:"product_code" gorm:"size:50"`
	ProductName string    `json:"product_name" gorm:"size:255"`
	
	Quantity   int     `json:"quantity"`
	UnitValue  float64 `json:"unit_value" gorm:"type:decimal(18,2)"`
	TotalValue float64 `json:"total_value" gorm:"type:decimal(18,2)"`
	
	LotID     *uuid.UUID `json:"lot_id" gorm:"type:uuid"`
	LotNumber string     `json:"lot_number" gorm:"size:100"`
	
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SampleItem) TableName() string {
	return "sample_items"
}

// CalculateTotal calculates total value for the item
func (i *SampleItem) CalculateTotal() {
	i.TotalValue = float64(i.Quantity) * i.UnitValue
}

// ShipmentStatus represents shipment status
type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "PENDING"
	ShipmentStatusShipped   ShipmentStatus = "SHIPPED"
	ShipmentStatusInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
	ShipmentStatusReturned  ShipmentStatus = "RETURNED"
)

// SampleShipment represents a shipment of samples to KOL
type SampleShipment struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShipmentNumber  string         `json:"shipment_number" gorm:"uniqueIndex;size:50"`
	SampleRequestID uuid.UUID      `json:"sample_request_id" gorm:"type:uuid;not null"`
	SampleRequest   *SampleRequest `json:"sample_request,omitempty" gorm:"foreignKey:SampleRequestID"`

	ShipmentDate   time.Time `json:"shipment_date"`
	Courier        string    `json:"courier" gorm:"size:100"`
	TrackingNumber string    `json:"tracking_number" gorm:"size:100"`

	RecipientName   string `json:"recipient_name" gorm:"size:200"`
	RecipientPhone  string `json:"recipient_phone" gorm:"size:50"`
	DeliveryAddress string `json:"delivery_address"`

	EstimatedDelivery *time.Time `json:"estimated_delivery"`
	ActualDelivery    *time.Time `json:"actual_delivery"`

	Status ShipmentStatus `json:"status" gorm:"size:50;default:PENDING"`

	DeliveryNotes   string `json:"delivery_notes"`
	ProofOfDelivery string `json:"proof_of_delivery"`

	CreatedBy *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (SampleShipment) TableName() string {
	return "sample_shipments"
}

// Ship marks shipment as shipped
func (s *SampleShipment) Ship(courier, trackingNumber string) {
	s.Status = ShipmentStatusShipped
	s.Courier = courier
	s.TrackingNumber = trackingNumber
	s.ShipmentDate = time.Now()
	s.UpdatedAt = time.Now()
}

// MarkDelivered marks shipment as delivered
func (s *SampleShipment) MarkDelivered() {
	now := time.Now()
	s.Status = ShipmentStatusDelivered
	s.ActualDelivery = &now
	s.UpdatedAt = now
}
