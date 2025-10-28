package producer

import (
	"context"
	"strconv"

	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/platform_common/pkg/logger"
	events_v1 "github.com/WithSoull/platform_common/pkg/proto/events/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func prettyDecoderUserDeleted(raw []byte) (string, bool) {
	msg := &events_v1.UserDeleted{}
	if err := proto.Unmarshal(raw, msg); err != nil {
		return "", false
	}
	j, err := PJ().Marshal(msg)
	if err != nil {
		return "", false
	}
	return string(j), true
}

func (s *userProducerService) ProduceUserDeleted(ctx context.Context, event model.UserDeletedEvent) error {
	var deletedAt *timestamppb.Timestamp
	if event.DeletedAt != nil {
		deletedAt = timestamppb.New(*event.DeletedAt)
	}

	msg := &events_v1.UserDeleted{
		UserId:    event.UserID,
		DeletedAt: deletedAt,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UserDeletedEvent", zap.Error(err))
		return err
	}

	err = s.userDeletedProducer.Send(ctx, []byte(strconv.FormatInt(event.UserID, 10)), payload, prettyDecoderUserDeleted)
	if err != nil {
		logger.Error(ctx, "failed to publish UserDeletedEvent", zap.Error(err))
		return err
	}

	return nil
}
