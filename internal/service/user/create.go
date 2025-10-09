package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/utils"
)

func (s *userService) Create(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) (int64, error) {
	// Input validation
	if userInfo.Name == "" {
		return 0, domainerrors.ErrNameRequired
	}
	if userInfo.Email == "" {
		return 0, domainerrors.ErrEmailRequired
	}
	if !utils.IsValidEmail(userInfo.Email) {
		return 0, domainerrors.ErrInvalidEmailFormat
	}
	if password == "" {
		return 0, domainerrors.ErrPasswordRequired
	}

	hashedPassword, err := s.validateAndHashPassword(password, passwordConfirm)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, &userInfo, hashedPassword)
	if err != nil {
		return 0, err
	}

	return id, nil
}
