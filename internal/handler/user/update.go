package user

import (
	"context"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	var name *string
	if req.GetName() != nil {
		name = &req.GetName().Value
	}

	var email *string
	if req.GetEmail() != nil {
		email = &req.GetEmail().Value
	}

	return &emptypb.Empty{}, s.service.Update(ctx, req.GetId(), name, email)
}
