package user

import (
	"context"
	"time"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/logger"
	"github.com/WithSoull/platform_common/pkg/sys/validate"
	"go.uber.org/zap"
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

	createdAt := time.Now()
	id, err := s.repo.Create(ctx, &userInfo, hashedPassword, createdAt)
	if err != nil {
		return 0, err
	}

	if err := s.userProducerService.ProduceUserCreated(ctx, model.UserCreatedEvent{
		UserID:    id,
		CreatedAt: &createdAt,
	}); err != nil {
		logger.Error(ctx, "failed to producer user.created", zap.Int64("userID", id), zap.Error(err))
	}

	return id, nil
}
