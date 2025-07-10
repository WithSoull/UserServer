package user

import (
	"context"

	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/repository/user/model"
	"github.com/WithSoull/AuthService/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) *service {
	return &service{
		repo: repo,
	}
}

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

	id, err := s.repo.Create(ctx, &userInfo, hashedPassword)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, id int64, name, email *string) error {
	if *name == "" {
			return status.Errorf(codes.InvalidArgument, "name is required")
	}
	if *email == "" {
			return status.Errorf(codes.InvalidArgument, "email is required")
	}
	if !utils.IsValidEmail(*email) {
			return status.Errorf(codes.InvalidArgument, "invalid email format")
	}
	return s.repo.Update(ctx, id, name, email)
}

func (s *service) UpdatePassword(ctx context.Context, id int64, password, passwordConfirm string) (error) {
	hashedPassword, err := s.hashPassword(password, passwordConfirm)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, id, hashedPassword)	
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) hashPassword(password, passwordConfirm string) (string, error) {
	if password != passwordConfirm {
			return "", status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
			return "", status.Errorf(codes.Internal, "failed to hash password")
	}

	return string(hashedPassword), nil
}
