package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *userService) UpdatePassword(ctx context.Context, id int64, password, passwordConfirm string) error {
	hashedPassword, err := s.validateAndHashPassword(password, passwordConfirm)
	if err != nil {
		return err
	}

	txErr := s.txManger.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdatePassword(ctx, id, hashedPassword); err != nil {
			return err
		}
		ip, ok := ipctx.ExtractIP(ctx)
		if !ok {
			log.Printf("[Service Layer] failed to extract IP from context for user (id=%d)", id)
		}
		return s.repo.LogPassword(ctx, id, ip)
	})
	if txErr != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(txErr)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to update user password: %v", txErr)
		}
		return grpcErr
	}

	return nil
}

// Validate password and hash it with grpc-code errorr
func (s *userService) validateAndHashPassword(password, passwordConfirm string) (string, error) {
	if password != passwordConfirm {
		return "", status.Error(codes.InvalidArgument, "passwords do not match")
	}

	if len(password) < 5 {
		return "", status.Error(codes.InvalidArgument, "password must be at least 5 characters long")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[Service Layer] failed to hash password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}
