package quotation

import (
	"context"
	"errors"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

var (
	ErrQuotationNotFound     = errors.New("quotation not found")
	ErrQuotationExpired      = errors.New("quotation has expired")
	ErrQuotationCannotSend   = errors.New("quotation cannot be sent")
	ErrQuotationCannotConvert = errors.New("quotation cannot be converted to order")
)

// CreateQuotationInput represents input for creating quotation
type CreateQuotationInput struct {
	CustomerID         uuid.UUID
	QuotationDate      time.Time
	ValidUntil         time.Time
	DiscountPercent    float64
	DiscountAmount     float64
	TaxPercent         float64
	Notes              string
	TermsAndConditions string
	Items              []QuotationItemInput
	CreatedBy          *uuid.UUID
}

// QuotationItemInput represents a line item input
type QuotationItemInput struct {
	ProductID       uuid.UUID
	ProductCode     string
	ProductName     string
	Quantity        float64
	UomID           *uuid.UUID
	UnitPrice       float64
	DiscountPercent float64
	TaxPercent      float64
	Notes           string
}

// CreateQuotationUseCase handles quotation creation
type CreateQuotationUseCase struct {
	quotationRepo repository.QuotationRepository
	customerRepo  repository.CustomerRepository
}

// NewCreateQuotationUseCase creates a new use case
func NewCreateQuotationUseCase(quotationRepo repository.QuotationRepository, customerRepo repository.CustomerRepository) *CreateQuotationUseCase {
	return &CreateQuotationUseCase{
		quotationRepo: quotationRepo,
		customerRepo:  customerRepo,
	}
}

// Execute creates a new quotation
func (uc *CreateQuotationUseCase) Execute(ctx context.Context, input *CreateQuotationInput) (*entity.Quotation, error) {
	// Verify customer exists
	_, err := uc.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return nil, err
	}

	// Generate quotation number
	number, err := uc.quotationRepo.GetNextQuotationNumber(ctx)
	if err != nil {
		return nil, err
	}

	quotation := &entity.Quotation{
		QuotationNumber:    number,
		CustomerID:         input.CustomerID,
		QuotationDate:      input.QuotationDate,
		ValidUntil:         input.ValidUntil,
		DiscountPercent:    input.DiscountPercent,
		DiscountAmount:     input.DiscountAmount,
		TaxPercent:         input.TaxPercent,
		Notes:              input.Notes,
		TermsAndConditions: input.TermsAndConditions,
		Status:             entity.QuotationStatusDraft,
		CreatedBy:          input.CreatedBy,
	}

	// Create line items
	for i, item := range input.Items {
		lineItem := entity.QuotationLineItem{
			LineNumber:      i + 1,
			ProductID:       item.ProductID,
			ProductCode:     item.ProductCode,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			UomID:           item.UomID,
			UnitPrice:       item.UnitPrice,
			DiscountPercent: item.DiscountPercent,
			TaxPercent:      item.TaxPercent,
			Notes:           item.Notes,
		}
		lineItem.CalculateLineTotal()
		quotation.LineItems = append(quotation.LineItems, lineItem)
	}

	// Calculate totals
	quotation.CalculateTotals()

	if err := uc.quotationRepo.Create(ctx, quotation); err != nil {
		return nil, err
	}

	return quotation, nil
}

// GetQuotationUseCase handles getting quotation
type GetQuotationUseCase struct {
	quotationRepo repository.QuotationRepository
}

// NewGetQuotationUseCase creates a new use case
func NewGetQuotationUseCase(repo repository.QuotationRepository) *GetQuotationUseCase {
	return &GetQuotationUseCase{quotationRepo: repo}
}

// Execute gets a quotation by ID
func (uc *GetQuotationUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Quotation, error) {
	return uc.quotationRepo.GetByID(ctx, id)
}

// ListQuotationsUseCase handles listing quotations
type ListQuotationsUseCase struct {
	quotationRepo repository.QuotationRepository
}

// NewListQuotationsUseCase creates a new use case
func NewListQuotationsUseCase(repo repository.QuotationRepository) *ListQuotationsUseCase {
	return &ListQuotationsUseCase{quotationRepo: repo}
}

// Execute lists quotations with filters
func (uc *ListQuotationsUseCase) Execute(ctx context.Context, filter *repository.QuotationFilter) ([]*entity.Quotation, int64, error) {
	return uc.quotationRepo.List(ctx, filter)
}

// SendQuotationUseCase handles sending quotation
type SendQuotationUseCase struct {
	quotationRepo repository.QuotationRepository
	eventPub      *event.Publisher
}

// NewSendQuotationUseCase creates a new use case
func NewSendQuotationUseCase(repo repository.QuotationRepository, eventPub *event.Publisher) *SendQuotationUseCase {
	return &SendQuotationUseCase{
		quotationRepo: repo,
		eventPub:      eventPub,
	}
}

// Execute sends a quotation
func (uc *SendQuotationUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Quotation, error) {
	quotation, err := uc.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !quotation.CanBeSent() {
		return nil, ErrQuotationCannotSend
	}

	quotation.Send()

	if err := uc.quotationRepo.Update(ctx, quotation); err != nil {
		return nil, err
	}

	// Publish event
	if uc.eventPub != nil {
		uc.eventPub.PublishQuotationSent(&event.QuotationSentEvent{
			QuotationID:     quotation.ID.String(),
			QuotationNumber: quotation.QuotationNumber,
			CustomerID:      quotation.CustomerID.String(),
			TotalAmount:     quotation.TotalAmount,
			ValidUntil:      quotation.ValidUntil.Format("2006-01-02"),
		})
	}

	return quotation, nil
}

// ConvertToOrderUseCase handles converting quotation to sales order
type ConvertToOrderUseCase struct {
	quotationRepo  repository.QuotationRepository
	salesOrderRepo repository.SalesOrderRepository
	customerRepo   repository.CustomerRepository
	eventPub       *event.Publisher
}

// NewConvertToOrderUseCase creates a new use case
func NewConvertToOrderUseCase(
	quotationRepo repository.QuotationRepository,
	salesOrderRepo repository.SalesOrderRepository,
	customerRepo repository.CustomerRepository,
	eventPub *event.Publisher,
) *ConvertToOrderUseCase {
	return &ConvertToOrderUseCase{
		quotationRepo:  quotationRepo,
		salesOrderRepo: salesOrderRepo,
		customerRepo:   customerRepo,
		eventPub:       eventPub,
	}
}

// ConvertToOrderInput represents input for converting quotation to order
type ConvertToOrderInput struct {
	QuotationID     uuid.UUID
	DeliveryDate    *time.Time
	DeliveryAddress string
	PaymentMethod   entity.PaymentMethod
	Notes           string
	CreatedBy       *uuid.UUID
}

// Execute converts quotation to sales order
func (uc *ConvertToOrderUseCase) Execute(ctx context.Context, input *ConvertToOrderInput) (*entity.SalesOrder, error) {
	quotation, err := uc.quotationRepo.GetByID(ctx, input.QuotationID)
	if err != nil {
		return nil, err
	}

	if !quotation.CanBeConverted() {
		return nil, ErrQuotationCannotConvert
	}

	if quotation.IsExpired() {
		return nil, ErrQuotationExpired
	}

	// Generate SO number
	soNumber, err := uc.salesOrderRepo.GetNextSONumber(ctx)
	if err != nil {
		return nil, err
	}

	// Create sales order from quotation
	salesOrder := &entity.SalesOrder{
		SONumber:        soNumber,
		CustomerID:      quotation.CustomerID,
		QuotationID:     &quotation.ID,
		SODate:          time.Now(),
		DeliveryDate:    input.DeliveryDate,
		DeliveryAddress: input.DeliveryAddress,
		Subtotal:        quotation.Subtotal,
		DiscountPercent: quotation.DiscountPercent,
		DiscountAmount:  quotation.DiscountAmount,
		TaxPercent:      quotation.TaxPercent,
		TaxAmount:       quotation.TaxAmount,
		TotalAmount:     quotation.TotalAmount,
		Status:          entity.SOStatusDraft,
		PaymentMethod:   input.PaymentMethod,
		PaymentStatus:   entity.PaymentStatusPending,
		Notes:           input.Notes,
		CreatedBy:       input.CreatedBy,
	}

	// Copy line items
	for i, item := range quotation.LineItems {
		soItem := entity.SOLineItem{
			LineNumber:      i + 1,
			ProductID:       item.ProductID,
			ProductCode:     item.ProductCode,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			UomID:           item.UomID,
			UnitPrice:       item.UnitPrice,
			DiscountPercent: item.DiscountPercent,
			DiscountAmount:  item.DiscountAmount,
			TaxPercent:      item.TaxPercent,
			TaxAmount:       item.TaxAmount,
			LineTotal:       item.LineTotal,
			Notes:           item.Notes,
		}
		salesOrder.LineItems = append(salesOrder.LineItems, soItem)
	}

	if err := uc.salesOrderRepo.Create(ctx, salesOrder); err != nil {
		return nil, err
	}

	// Mark quotation as converted
	quotation.MarkConverted()
	if err := uc.quotationRepo.Update(ctx, quotation); err != nil {
		return nil, err
	}

	// Publish order created event
	if uc.eventPub != nil {
		items := make([]event.OrderLineItem, len(salesOrder.LineItems))
		for i, item := range salesOrder.LineItems {
			items[i] = event.OrderLineItem{
				ProductID:   item.ProductID.String(),
				ProductCode: item.ProductCode,
				ProductName: item.ProductName,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			}
		}
		uc.eventPub.PublishOrderCreated(&event.OrderCreatedEvent{
			SOID:        salesOrder.ID.String(),
			SONumber:    salesOrder.SONumber,
			CustomerID:  salesOrder.CustomerID.String(),
			TotalAmount: salesOrder.TotalAmount,
			Items:       items,
		})
	}

	return salesOrder, nil
}
