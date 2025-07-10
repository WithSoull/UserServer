package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/repository/user/conventer"
	"github.com/WithSoull/AuthService/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) Create(ctx context.Context, userInfo model.UserInfo, password, passwordConfirm string) (int64, error) {
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
	if userInfo.Role < 0 && userInfo.Role > 1 { // Assuming role is an enum with positive values
			return 0, status.Errorf(codes.InvalidArgument, "invalid role")
	}
	
	hashedPassword, err := s.hashPassword(password, passwordConfirm)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.Create(ctx, conventer.FromModelToRepoUserInfo(&userInfo), hashedPassword)
	if err != nil {
		return 0, err
	}

	return id, nil
}
