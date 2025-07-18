package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/contextx/ipctx"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	ctx = ipctx.InjectIp(ctx)
	return &emptypb.Empty{}, s.service.UpdatePassword(ctx, req.GetId(), req.GetPassword(), req.GetPasswordConfirm())
}
