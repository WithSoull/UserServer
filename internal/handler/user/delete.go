package user

import (
	"context"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	ctx, err := h.verifyToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, h.service.Delete(ctx, req.GetId())
}
