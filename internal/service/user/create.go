package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
)

func (s *userService) Create(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) (int64, error) {
	// UserInfo Validation
	if err := validate.Validate(
		ctx,
		validator.ValidateNotEmptyString(userInfo.Name, domainerrors.ErrNameRequired),
		validator.ValidateNotEmptyString(userInfo.Email, domainerrors.ErrEmailRequired),
		validator.ValidateEmailFromat(userInfo.Email),
	); err != nil {
		return 0, err
	}

	// Password Validation + Hashing
	hashedPassword, err := s.validateAndHashPassword(ctx, password, passwordConfirm)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, &userInfo, hashedPassword)
	if err != nil {
		return 0, err
	}

	return id, nil
}
