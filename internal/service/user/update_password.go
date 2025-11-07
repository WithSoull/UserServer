package user

import (
	"context"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/WithSoull/UserServer/internal/validator"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/contextx/ipctx"
	"github.com/WithSoull/platform_common/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) UpdatePassword(ctx context.Context, password, passwordConfirm string) error {
	senderID, ok := claimsctx.ExtractUserID(ctx)
	if !ok {
		return domainerrors.ErrFailedToVerify
	}

	// Input Validation + Hashing
	hashedPassword, err := s.validateAndHashPassword(ctx, password, passwordConfirm)
	if err != nil {
		return err
	}

	txErr := s.txManger.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.repo.UpdatePassword(ctx, senderID, hashedPassword); err != nil {
			return err
		}
		ip, ok := ipctx.ExtractIP(ctx)
		if !ok {
			logger.Error(claimsctx.InjectUserID(ctx, senderID), "Failed to extract IP from context to user")
		}
		return s.repo.LogPassword(ctx, senderID, ip)
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
