package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"golang.org/x/crypto/bcrypt"
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

	return txErr
}

// Validate password and hash it with grpc-code errorr
func (s *userService) validateAndHashPassword(password, passwordConfirm string) (string, error) {
	// Input Validation
	if password == "" {
		return "", domainerrors.ErrPasswordRequired
	}
	if password != passwordConfirm {
		return "", domainerrors.ErrPasswordMismatch
	}

	if len(password) < 5 {
		return "", domainerrors.ErrPasswordTooShort
	}

	return s.hashPassword(password)
}

func (s *userService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}
