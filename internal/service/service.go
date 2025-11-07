package service

import (
	"context"

	"github.com/WithSoull/UserServer/internal/model"
)

type UserService interface {
	Create(context.Context, model.UserInfo, string, string) (int64, error)
	Get(context.Context, int64) (*model.User, error)
	Update(context.Context, *string, *string) error
	UpdatePassword(context.Context, string, string) error
	Delete(context.Context) error
	ValidateCredentials(context.Context, string, string) (bool, int64)
}

type UserProducerService interface {
	ProduceUserCreated(ctx context.Context, event model.UserCreatedEvent) error
	ProduceUserDeleted(ctx context.Context, event model.UserDeletedEvent) error
}
