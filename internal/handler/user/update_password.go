package user

import (
	"context"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	ctx, err := h.verifyToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	ctx = ipctx.InjectIp(ctx)
	return &emptypb.Empty{}, h.service.UpdatePassword(ctx, req.GetPassword(), req.GetPasswordConfirm())
}
