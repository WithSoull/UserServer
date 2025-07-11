package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Update(ctx context.Context, id int64, name, email *string) error {
	if name != nil && *name == "" {
			return status.Errorf(codes.InvalidArgument, "name is required")
	}
	if email != nil && *email == "" {
			return status.Errorf(codes.InvalidArgument, "email is required")
	}
	if email != nil && !utils.IsValidEmail(*email) {
			return status.Errorf(codes.InvalidArgument, "invalid email format")
	}
	return s.repo.Update(ctx, id, name, email)
}
