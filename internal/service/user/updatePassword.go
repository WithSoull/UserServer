package user

import (
	"context"
	"log"

	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) UpdatePassword(ctx context.Context, id int64, password, passwordConfirm string) error {
	// Input Validation + Hashing
	hashedPassword, err := s.validateAndHashPassword(ctx, password, passwordConfirm)
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

func (s *userService) validateAndHashPassword(ctx context.Context, password, passwordConfirm string) (string, error) {
	if err := validator.ValidatePassword(ctx, password, passwordConfirm); err != nil {
		return "", err
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
