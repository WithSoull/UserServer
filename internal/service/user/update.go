package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain_errors"
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
	err := s.repo.Update(ctx, id, name, email)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to update user: %v", err)
		}
		return grpcErr
	}
	return nil
}
