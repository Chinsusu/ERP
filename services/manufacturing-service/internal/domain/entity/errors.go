package entity

// Error definitions for manufacturing domain
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// Common domain errors
var (
	ErrBOMNotFound             = &DomainError{Code: "BOM_NOT_FOUND", Message: "BOM not found"}
	ErrBOMNotDraft             = &DomainError{Code: "BOM_NOT_DRAFT", Message: "BOM is not in draft status"}
	ErrBOMNotPendingApproval   = &DomainError{Code: "BOM_NOT_PENDING", Message: "BOM is not pending approval"}
	ErrBOMAlreadyApproved      = &DomainError{Code: "BOM_ALREADY_APPROVED", Message: "BOM is already approved"}
	ErrInvalidEncryptionKey    = &DomainError{Code: "INVALID_KEY", Message: "Invalid encryption key"}
	
	ErrWONotFound              = &DomainError{Code: "WO_NOT_FOUND", Message: "Work order not found"}
	ErrWOCannotRelease         = &DomainError{Code: "WO_CANNOT_RELEASE", Message: "Work order cannot be released"}
	ErrWOCannotStart           = &DomainError{Code: "WO_CANNOT_START", Message: "Work order cannot be started"}
	ErrWOCannotComplete        = &DomainError{Code: "WO_CANNOT_COMPLETE", Message: "Work order cannot be completed"}
	ErrWOCannotCancel          = &DomainError{Code: "WO_CANNOT_CANCEL", Message: "Work order cannot be cancelled"}
	ErrMaterialNotIssued       = &DomainError{Code: "MATERIAL_NOT_ISSUED", Message: "Material not fully issued"}
	
	ErrQCInspectionNotFound    = &DomainError{Code: "QC_NOT_FOUND", Message: "QC inspection not found"}
	ErrQCAlreadyApproved       = &DomainError{Code: "QC_ALREADY_APPROVED", Message: "QC inspection already approved"}
	ErrQCCheckpointNotFound    = &DomainError{Code: "QC_CHECKPOINT_NOT_FOUND", Message: "QC checkpoint not found"}
	
	ErrNCRNotFound             = &DomainError{Code: "NCR_NOT_FOUND", Message: "NCR not found"}
	ErrNCRAlreadyClosed        = &DomainError{Code: "NCR_ALREADY_CLOSED", Message: "NCR is already closed"}
	
	ErrTraceNotFound           = &DomainError{Code: "TRACE_NOT_FOUND", Message: "Traceability record not found"}
	ErrLotNotFound             = &DomainError{Code: "LOT_NOT_FOUND", Message: "Lot not found"}
)
