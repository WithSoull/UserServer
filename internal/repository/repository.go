package repository

import (
	"context"

	"github.com/WithSoull/AuthService/internal/repository/user/model"
)


type UserRepository interface {
	Create(ctx context.Context, user *model.User, hashedPassword string) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, name, email *string) (error)
	UpdatePassword(ctx context.Context, id int64, hashedPassword string) error
	Delete(ctx context.Context, id int64) (error)
}
