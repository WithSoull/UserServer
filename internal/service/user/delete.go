package user

import (
	"context"
	"time"

	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
	"go.uber.org/zap"
)

func (s *userService) Delete(ctx context.Context, id int64) error {
	if err := s.checkUserPermission(ctx, id); err != nil {
		return err
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	if err := s.userProducerService.ProduceUserDeleted(ctx, model.UserDeletedEvent{
		UserID:    id,
		DeletedAt: &now,
	}); err != nil {
		logger.Error(ctx, "failed to producer user.deleted", zap.Int64("userID", id), zap.Error(err))
	}

	return nil
}
