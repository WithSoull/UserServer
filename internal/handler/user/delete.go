package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) Delete(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	ctx, err := h.verifyToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, h.service.Delete(ctx)
}
