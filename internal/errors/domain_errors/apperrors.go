package domainerrors

import (
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
)

var (
	// Resource errors (NotFound)
	ErrUserNotFound = sys.NewCommonError("user not found", codes.NotFound)

	// Conflict errors (AlreadyExists)
	ErrEmailAlreadyExists = sys.NewCommonError("email already exists", codes.AlreadyExists)

	// Validation errors - General (InvalidArgument)
	ErrInvalidInput = sys.NewCommonError("invalid input", codes.InvalidArgument)

	// Validation errors - Required fields (InvalidArgument)
	ErrNameRequired     = sys.NewCommonError("name is required", codes.InvalidArgument)
	ErrEmailRequired    = sys.NewCommonError("email is required", codes.InvalidArgument)
	ErrPasswordRequired = sys.NewCommonError("password is required", codes.InvalidArgument)

	// Validation errors - Format (InvalidArgument)
	ErrInvalidEmailFormat = sys.NewCommonError("invalid email format", codes.InvalidArgument)

	// Validation errors - Password (InvalidArgument)
	ErrPasswordMismatch = sys.NewCommonError("passwords do not match", codes.InvalidArgument)
	ErrPasswordTooShort = sys.NewCommonError("password must be at least 5 characters long", codes.InvalidArgument)

	// Internal errors
	ErrInternal = sys.NewCommonError("internal error", codes.Internal)
)
