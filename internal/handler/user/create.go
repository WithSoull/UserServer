package user

import (
	"context"

	conventer "github.com/WithSoull/AuthService/internal/conventer/user"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
)

func (s *handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userID, err := s.service.Create(ctx, conventer.FromProtoToModelUserInfo(req.GetUserInfo()), req.GetPassword(), req.GetPasswordConfirm())
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}
