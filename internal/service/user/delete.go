package user

import (
	"context"
	"time"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (s *userService) Delete(ctx context.Context) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	err := s.repo.Delete(ctx, senderID)
	if err != nil {
		return err
	}

	now := time.Now()
	if err := s.userProducerService.ProduceUserDeleted(ctx, model.UserDeletedEvent{
		UserID:    senderID,
		DeletedAt: &now,
	}); err != nil {
		logger.Error(ctx, "failed to producer user.deleted", zap.Int64("userID", senderID), zap.Error(err))
	}

	return nil
}
