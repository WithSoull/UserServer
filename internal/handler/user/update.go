package user

import (
	"context"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	ctx, err := h.verifyToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	var name *string
	if req.GetName() != nil {
		name = &req.GetName().Value
	}

	var email *string
	if req.GetEmail() != nil {
		email = &req.GetEmail().Value
	}

	return &emptypb.Empty{}, h.service.Update(ctx, name, email)
}
