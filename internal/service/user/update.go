package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func (s *userService) Update(ctx context.Context, id int64, name, email *string) error {
	// Input validation
	if err := validate.Validate(
		ctx,
		validator.ValidateNotEmptyPointerToString(name, domainerrors.ErrNameRequired),
		validator.ValidateNotEmptyPointerToString(email, domainerrors.ErrEmailRequired),
		validator.ValidateEmailFromatPointer(email),
	); err != nil {
		return err
	}

	return s.repo.Update(ctx, id, name, email)
}
