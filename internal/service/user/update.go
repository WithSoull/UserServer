package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/utils"
)

func (s *userService) Update(ctx context.Context, id int64, name, email *string) error {
	// Input validation
	if name != nil && *name == "" {
		return domainerrors.ErrNameRequired
	}
	if email != nil && *email == "" {
		return domainerrors.ErrEmailRequired
	}
	if email != nil && !utils.IsValidEmail(*email) {
		return domainerrors.ErrInvalidEmailFormat
	}

	return s.repo.Update(ctx, id, name, email)
}
