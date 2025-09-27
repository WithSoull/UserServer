package user

import (
	"context"

	conventer "github.com/WithSoull/UserServer/internal/conventer/user"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
)

func (s *handler) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: conventer.FromModelToProtoUser(*user),
	}, nil
}
