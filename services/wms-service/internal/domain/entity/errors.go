package entity

import "errors"

// Domain errors
var (
	ErrNotFound          = errors.New("not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrLotExpired        = errors.New("lot is expired")
	ErrLotNotAvailable   = errors.New("lot is not available")
	ErrQCNotPassed       = errors.New("QC not passed")
	ErrAlreadyCompleted  = errors.New("already completed")
	ErrAlreadyCancelled  = errors.New("already cancelled")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidQuantity   = errors.New("invalid quantity")
	ErrReservationFailed = errors.New("reservation failed")
	ErrColdStorageAlert  = errors.New("cold storage temperature out of range")
	ErrPendingItems      = errors.New("pending items exist")
)
