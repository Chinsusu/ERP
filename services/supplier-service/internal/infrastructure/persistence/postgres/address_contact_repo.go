package postgres

import (
	"context"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *entity.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Address, error) {
	var address entity.Address
	err := r.db.WithContext(ctx).First(&address, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Address, error) {
	var addresses []*entity.Address
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("is_primary DESC, created_at").
		Find(&addresses).Error
	return addresses, err
}

func (r *addressRepository) Update(ctx context.Context, address *entity.Address) error {
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Address{}, "id = ?", id).Error
}

func (r *addressRepository) SetPrimary(ctx context.Context, supplierID, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all primary for this supplier
		if err := tx.Model(&entity.Address{}).
			Where("supplier_id = ?", supplierID).
			Update("is_primary", false).Error; err != nil {
			return err
		}
		// Set the specified address as primary
		return tx.Model(&entity.Address{}).
			Where("id = ?", addressID).
			Update("is_primary", true).Error
	})
}

type contactRepository struct {
	db *gorm.DB
}

// NewContactRepository creates a new contact repository
func NewContactRepository(db *gorm.DB) repository.ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) Create(ctx context.Context, contact *entity.Contact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

func (r *contactRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Contact, error) {
	var contact entity.Contact
	err := r.db.WithContext(ctx).First(&contact, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *contactRepository) GetBySupplierID(ctx context.Context, supplierID uuid.UUID) ([]*entity.Contact, error) {
	var contacts []*entity.Contact
	err := r.db.WithContext(ctx).
		Where("supplier_id = ?", supplierID).
		Order("is_primary DESC, created_at").
		Find(&contacts).Error
	return contacts, err
}

func (r *contactRepository) Update(ctx context.Context, contact *entity.Contact) error {
	return r.db.WithContext(ctx).Save(contact).Error
}

func (r *contactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Contact{}, "id = ?", id).Error
}

func (r *contactRepository) SetPrimary(ctx context.Context, supplierID, contactID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all primary for this supplier
		if err := tx.Model(&entity.Contact{}).
			Where("supplier_id = ?", supplierID).
			Update("is_primary", false).Error; err != nil {
			return err
		}
		// Set the specified contact as primary
		return tx.Model(&entity.Contact{}).
			Where("id = ?", contactID).
			Update("is_primary", true).Error
	})
}
