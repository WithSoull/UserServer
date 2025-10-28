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

func prettyDecoderUserCreated(raw []byte) (string, bool) {
	msg := &events_v1.UserCreated{}
	if err := proto.Unmarshal(raw, msg); err != nil {
		return "", false
	}
	j, err := PJ().Marshal(msg)
	if err != nil {
		return "", false
	}
	return string(j), true
}

func (s *userProducerService) ProduceUserCreated(ctx context.Context, event model.UserCreatedEvent) error {
	var createdAt *timestamppb.Timestamp
	if event.CreatedAt != nil {
		createdAt = timestamppb.New(*event.CreatedAt)
	}

	msg := &events_v1.UserCreated{
		UserId:    event.UserID,
		CreatedAt: createdAt,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal UserCreatedEvent", zap.Error(err))
		return err
	}

	err = s.userCreatedProducer.Send(ctx, []byte(strconv.FormatInt(event.UserID, 10)), payload, prettyDecoderUserCreated)
	if err != nil {
		logger.Error(ctx, "failed to publish UserCreatedEvent", zap.Error(err))
		return err
	}

	return nil
}
