package user

import (
	"context"

	"github.com/WithSoull/UserServer/internal/model"
)

func (s *userService) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
