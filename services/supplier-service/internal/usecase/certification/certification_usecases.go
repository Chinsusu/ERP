package certification

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// AddCertificationRequest represents the request to add a certification
type AddCertificationRequest struct {
	SupplierID   uuid.UUID `json:"-"`
	Type         string    `json:"certification_type" validate:"required,oneof=GMP ISO9001 ISO22716 ORGANIC ECOCERT HALAL COSMOS OTHER"`
	CertNumber   string    `json:"certificate_number" validate:"required"`
	IssuingBody  string    `json:"issuing_body" validate:"required"`
	IssueDate    string    `json:"issue_date" validate:"required"`
	ExpiryDate   string    `json:"expiry_date" validate:"required"`
	DocumentURL  string    `json:"document_url"`
	Notes        string    `json:"notes"`
}

// AddCertificationUseCase handles adding a certification to a supplier
type AddCertificationUseCase struct {
	certRepo     repository.CertificationRepository
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewAddCertificationUseCase creates a new AddCertificationUseCase
func NewAddCertificationUseCase(
	certRepo repository.CertificationRepository,
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *AddCertificationUseCase {
	return &AddCertificationUseCase{
		certRepo:     certRepo,
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute adds a certification
func (uc *AddCertificationUseCase) Execute(ctx context.Context, req *AddCertificationRequest) (*entity.Certification, error) {
	// Get supplier to validate it exists
	supplier, err := uc.supplierRepo.GetByID(ctx, req.SupplierID)
	if err != nil {
		return nil, err
	}

	// Parse dates
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		return nil, err
	}
	expiryDate, err := time.Parse("2006-01-02", req.ExpiryDate)
	if err != nil {
		return nil, err
	}

	cert := &entity.Certification{
		ID:          uuid.New(),
		SupplierID:  req.SupplierID,
		Type:        entity.CertificationType(req.Type),
		CertNumber:  req.CertNumber,
		IssuingBody: req.IssuingBody,
		IssueDate:   issueDate,
		ExpiryDate:  expiryDate,
		DocumentURL: req.DocumentURL,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Status is auto-calculated
	cert.UpdateStatus()

	if err := uc.certRepo.Create(ctx, cert); err != nil {
		return nil, err
	}

	// Publish event
	uc.eventPub.PublishCertificationAdded(ctx, &event.CertificationEvent{
		SupplierID:        supplier.ID.String(),
		SupplierName:      supplier.Name,
		CertificationType: string(cert.Type),
		CertNumber:        cert.CertNumber,
		ExpiryDate:        cert.ExpiryDate,
	})

	return cert, nil
}

// GetCertificationsUseCase handles getting certifications for a supplier
type GetCertificationsUseCase struct {
	certRepo repository.CertificationRepository
}

// NewGetCertificationsUseCase creates a new GetCertificationsUseCase
func NewGetCertificationsUseCase(certRepo repository.CertificationRepository) *GetCertificationsUseCase {
	return &GetCertificationsUseCase{certRepo: certRepo}
}

// Execute gets certifications for a supplier
func (uc *GetCertificationsUseCase) Execute(ctx context.Context, supplierID uuid.UUID) ([]*entity.Certification, error) {
	return uc.certRepo.GetBySupplierID(ctx, supplierID)
}

// GetExpiringCertificationsUseCase handles getting expiring certifications
type GetExpiringCertificationsUseCase struct {
	certRepo repository.CertificationRepository
}

// NewGetExpiringCertificationsUseCase creates a new GetExpiringCertificationsUseCase
func NewGetExpiringCertificationsUseCase(certRepo repository.CertificationRepository) *GetExpiringCertificationsUseCase {
	return &GetExpiringCertificationsUseCase{certRepo: certRepo}
}

// Execute gets certifications expiring within N days
func (uc *GetExpiringCertificationsUseCase) Execute(ctx context.Context, days int) ([]*entity.Certification, error) {
	if days <= 0 {
		days = 90 // Default to 90 days
	}
	return uc.certRepo.GetExpiring(ctx, days)
}

// CheckCertificationExpiryUseCase handles the scheduled job to check and update certification statuses
type CheckCertificationExpiryUseCase struct {
	certRepo     repository.CertificationRepository
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewCheckCertificationExpiryUseCase creates a new CheckCertificationExpiryUseCase
func NewCheckCertificationExpiryUseCase(
	certRepo repository.CertificationRepository,
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *CheckCertificationExpiryUseCase {
	return &CheckCertificationExpiryUseCase{
		certRepo:     certRepo,
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute runs the daily certificate expiry check
func (uc *CheckCertificationExpiryUseCase) Execute(ctx context.Context) error {
	now := time.Now()
	
	// Update expired certificates
	_, err := uc.certRepo.UpdateExpiredStatuses(ctx, now)
	if err != nil {
		return err
	}

	// Update expiring soon certificates (within 90 days)
	expiringCutoff := now.AddDate(0, 0, 90)
	_, err = uc.certRepo.UpdateExpiringStatuses(ctx, expiringCutoff)
	if err != nil {
		return err
	}

	// Get expiring certificates and publish events
	expiringCerts, err := uc.certRepo.GetExpiring(ctx, 90)
	if err != nil {
		return err
	}

	for _, cert := range expiringCerts {
		daysUntilExpiry := cert.CalculateDaysUntilExpiry()
		
		// Publish event for 90, 30, 7 day thresholds
		if daysUntilExpiry == 90 || daysUntilExpiry == 30 || daysUntilExpiry == 7 {
			supplier, _ := uc.supplierRepo.GetByID(ctx, cert.SupplierID)
			supplierName := ""
			if supplier != nil {
				supplierName = supplier.Name
			}
			
			uc.eventPub.PublishCertificationExpiring(ctx, &event.CertificationEvent{
				SupplierID:        cert.SupplierID.String(),
				SupplierName:      supplierName,
				CertificationType: string(cert.Type),
				CertNumber:        cert.CertNumber,
				ExpiryDate:        cert.ExpiryDate,
				DaysUntilExpiry:   daysUntilExpiry,
			})
		}
	}

	return nil
}
