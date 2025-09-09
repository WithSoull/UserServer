package user

import (
	"context"

	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	ctx = ipctx.InjectIp(ctx)
	return &emptypb.Empty{}, s.service.UpdatePassword(ctx, req.GetId(), req.GetPassword(), req.GetPasswordConfirm())
}
