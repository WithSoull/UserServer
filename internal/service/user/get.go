package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/repository/user/conventer"
)

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return conventer.FromRepoToModelUser(user), nil
}
