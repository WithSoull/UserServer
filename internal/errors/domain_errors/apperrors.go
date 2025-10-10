package domainerrors

import (
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

var (
	// Resource errors (NotFound)
	ErrUserNotFound = sys.NewCommonError("user not found", codes.NotFound)

	// Conflict errors (AlreadyExists)
	ErrEmailAlreadyExists = sys.NewCommonError("email already exists", codes.AlreadyExists)

	// Validation errors - General (InvalidArgument)
	ErrInvalidInput = validate.NewValidationErrors("invalid input")

	// Validation errors - Required fields (InvalidArgument)
	ErrNameRequired      = validate.NewValidationErrors("name is required")
	ErrEmailRequired     = validate.NewValidationErrors("email is required")
	ErrPasswordRequired  = validate.NewValidationErrors("password is required")
	ErrNoChangesProvided = validate.NewValidationErrors("no changes provided")

	// Validation errors - Format (InvalidArgument)
	ErrInvalidEmailFormat = validate.NewValidationErrors("invalid email format")

	// Validation errors - Password (InvalidArgument)
	ErrPasswordMismatch = validate.NewValidationErrors("passwords do not match")
	ErrPasswordTooShort = validate.NewValidationErrors("password must be at least 5 characters long")

	// Internal errors
	ErrInternal = sys.NewCommonError("internal error", codes.Internal)
)
