package entity

import "errors"

// Common errors
var (
	ErrInvalidPRStatus = errors.New("invalid PR status for this action")
	ErrInvalidPOStatus = errors.New("invalid PO status for this action")
	ErrPRNotApproved   = errors.New("PR must be approved before conversion")
)
