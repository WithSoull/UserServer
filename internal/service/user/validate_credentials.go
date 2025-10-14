package user

import (
	"context"
	"errors"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) ValidateCredentials(ctx context.Context, email, password string) (bool, int64) {
	id, storedHash, err := s.repo.GetUserCredentials(ctx, email)
	if err != nil {
		if !errors.Is(err, domainerrors.ErrUserNotFound) {
			logger.Error(claimsctx.InjectUserEmail(ctx, email), "failed to get user credentials", zap.Error(err))
		}
		return false, 0
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false, 0
	}

	return true, id
}
