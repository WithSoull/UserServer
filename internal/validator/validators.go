package validator

import (
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
	"golang.org/x/net/context"
)

func ValidatePassword(ctx context.Context, password, passwordConfirm string) error {
	return validate.Validate(
		ctx,
		ValidateNotEmptyString(password, domainerrors.ErrPasswordRequired),
		ValidatePasswordTooShort(password),
		ValidatePasswordMismatch(password, passwordConfirm),
	)
}
