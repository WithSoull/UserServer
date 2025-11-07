package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func (s *userService) Update(ctx context.Context, name, email *string) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	// Input validation
	if err := validate.Validate(
		ctx,
		validator.ValidateNotEmptyPointerToString(name, domainerrors.ErrNameRequired),
		validator.ValidateNotEmptyPointerToString(email, domainerrors.ErrEmailRequired),
		validator.ValidateEmailFromatPointer(email),
	); err != nil {
		return err
	}

	return s.repo.Update(ctx, senderID, name, email)
}
