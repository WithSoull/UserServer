package user

import (
	"github.com/WithSoull/AuthService/internal/repository"
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
