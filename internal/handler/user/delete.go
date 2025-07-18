package user

import (
	"context"

	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.service.Delete(ctx, req.GetId())
}
