package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *entity.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).
		Preload("CustomerGroup").
		Preload("Addresses").
		Preload("Contacts").
		First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByCode(ctx context.Context, code string) (*entity.Customer, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).
		Preload("CustomerGroup").
		First(&customer, "customer_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	customer.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Customer{}, "id = ?", id).Error
}

func (r *customerRepository) List(ctx context.Context, filter *repository.CustomerFilter) ([]*entity.Customer, int64, error) {
	var customers []*entity.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Customer{})

	// Apply filters
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR customer_code ILIKE ? OR email ILIKE ?", search, search, search)
	}
	if filter.CustomerType != "" {
		query = query.Where("customer_type = ?", filter.CustomerType)
	}
	if filter.CustomerGroupID != nil {
		query = query.Where("customer_group_id = ?", filter.CustomerGroupID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Get results
	err := query.Preload("CustomerGroup").Order("created_at DESC").Find(&customers).Error
	return customers, total, err
}

func (r *customerRepository) GetNextCustomerCode(ctx context.Context) (string, error) {
	var count int64
	r.db.WithContext(ctx).Model(&entity.Customer{}).Count(&count)
	return fmt.Sprintf("CUST-%04d", count+1), nil
}

func (r *customerRepository) UpdateBalance(ctx context.Context, customerID uuid.UUID, amount float64) error {
	return r.db.WithContext(ctx).
		Model(&entity.Customer{}).
		Where("id = ?", customerID).
		Update("current_balance", gorm.Expr("current_balance + ?", amount)).
		Error
}

func (r *customerRepository) GetAvailableCredit(ctx context.Context, customerID uuid.UUID) (float64, error) {
	var customer entity.Customer
	err := r.db.WithContext(ctx).Select("credit_limit", "current_balance").First(&customer, "id = ?", customerID).Error
	if err != nil {
		return 0, err
	}
	return customer.CreditLimit - customer.CurrentBalance, nil
}

// Address operations
func (r *customerRepository) CreateAddress(ctx context.Context, address *entity.CustomerAddress) error {
	return r.db.WithContext(ctx).Create(address).Error
}

func (r *customerRepository) GetAddresses(ctx context.Context, customerID uuid.UUID) ([]*entity.CustomerAddress, error) {
	var addresses []*entity.CustomerAddress
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&addresses).Error
	return addresses, err
}

func (r *customerRepository) UpdateAddress(ctx context.Context, address *entity.CustomerAddress) error {
	address.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(address).Error
}

func (r *customerRepository) DeleteAddress(ctx context.Context, addressID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CustomerAddress{}, "id = ?", addressID).Error
}

func (r *customerRepository) GetDefaultAddress(ctx context.Context, customerID uuid.UUID, addressType entity.AddressType) (*entity.CustomerAddress, error) {
	var address entity.CustomerAddress
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND (address_type = ? OR address_type = 'BOTH') AND is_default = true", customerID, addressType).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Contact operations
func (r *customerRepository) CreateContact(ctx context.Context, contact *entity.CustomerContact) error {
	return r.db.WithContext(ctx).Create(contact).Error
}

func (r *customerRepository) GetContacts(ctx context.Context, customerID uuid.UUID) ([]*entity.CustomerContact, error) {
	var contacts []*entity.CustomerContact
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&contacts).Error
	return contacts, err
}

func (r *customerRepository) UpdateContact(ctx context.Context, contact *entity.CustomerContact) error {
	contact.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(contact).Error
}

func (r *customerRepository) DeleteContact(ctx context.Context, contactID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CustomerContact{}, "id = ?", contactID).Error
}

func (r *customerRepository) GetPrimaryContact(ctx context.Context, customerID uuid.UUID) (*entity.CustomerContact, error) {
	var contact entity.CustomerContact
	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND is_primary = true", customerID).
		First(&contact).Error
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// CustomerGroupRepository implementation
type customerGroupRepository struct {
	db *gorm.DB
}

// NewCustomerGroupRepository creates a new customer group repository
func NewCustomerGroupRepository(db *gorm.DB) repository.CustomerGroupRepository {
	return &customerGroupRepository{db: db}
}

func (r *customerGroupRepository) Create(ctx context.Context, group *entity.CustomerGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *customerGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CustomerGroup, error) {
	var group entity.CustomerGroup
	err := r.db.WithContext(ctx).First(&group, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *customerGroupRepository) GetByCode(ctx context.Context, code string) (*entity.CustomerGroup, error) {
	var group entity.CustomerGroup
	err := r.db.WithContext(ctx).First(&group, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *customerGroupRepository) Update(ctx context.Context, group *entity.CustomerGroup) error {
	group.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(group).Error
}

func (r *customerGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CustomerGroup{}, "id = ?", id).Error
}

func (r *customerGroupRepository) List(ctx context.Context, activeOnly bool) ([]*entity.CustomerGroup, error) {
	var groups []*entity.CustomerGroup
	query := r.db.WithContext(ctx)
	if activeOnly {
		query = query.Where("is_active = true")
	}
	err := query.Order("code").Find(&groups).Error
	return groups, err
}
