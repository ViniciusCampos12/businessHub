package fails

import "errors"

var (
	ErrCompanyAlreadyExists          = errors.New("company already exists")
	ErrCompanyNotFound               = errors.New("company not found")
	ErrCompanyNotUpdated             = errors.New("company not updated")
	ErrCompanyDocumentIsAlreadyInUse = errors.New("company document is already in use")
	ErrInsufficientPWDQuota          = errors.New("insufficient quota")
	ErrDbDeleteFailed                = errors.New("mongodb deleted failed")
	ErrDbUpdateFailed                = errors.New("mongodb update failed")
	ErrValidation                    = errors.New("validation failed")
)
