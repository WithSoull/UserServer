package repository

import (
	"context"

	"github.com/WithSoull/UserServer/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.UserInfo, hashedPassword string) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	UpdatePassword(ctx context.Context, id int64, hashedPassword string) error
	LogPassword(ctx context.Context, id int64, ip_address string) error
	Delete(ctx context.Context, id int64) error
}
