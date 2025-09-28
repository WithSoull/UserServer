package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userService) Create(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) (int64, error) {
	// Input validation
	if userInfo.Name == "" {
		return 0, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if userInfo.Email == "" {
		return 0, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if !utils.IsValidEmail(userInfo.Email) {
		return 0, status.Errorf(codes.InvalidArgument, "invalid email format")
	}
	if password == "" {
		return 0, status.Errorf(codes.InvalidArgument, "password is required")
	}

	hashedPassword, err := s.validateAndHashPassword(password, passwordConfirm)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, &userInfo, hashedPassword)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to create user: %v", err)
		}
		return 0, grpcErr
	}

	return id, nil
}
