package user

import (
	"context"
	"log"

	"github.com/WithSoull/AuthService/internal/contextx/ipctx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *service) UpdatePassword(ctx context.Context, id int64, password, passwordConfirm string) error {
	hashedPassword, err := s.hashPassword(password, passwordConfirm)
	if err != nil {
		return err
	}
	return s.txManger.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdatePassword(ctx, id, hashedPassword); err != nil {
			return err
		}
		ip, ok := ipctx.ExtractIP(ctx)
		if !ok {
			log.Printf("failed to take ip from user(id = %d)", id)
		}
		return s.repo.LogPassword(ctx, id, ip)
	})
}

// Hash password and also check password is eq passwordConfirm
func (s *service) hashPassword(password, passwordConfirm string) (string, error) {
	if password != passwordConfirm {
		return "", status.Error(codes.InvalidArgument, "passwords does not match")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to hash password")
	}

	return string(hashedPassword), nil
}
