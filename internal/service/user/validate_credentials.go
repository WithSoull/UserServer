package user

import (
	"context"
	"errors"
	"log"

	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *userService) ValidateCredentials(ctx context.Context, email, password string) (bool, int64) {
	id, storedHash, err := s.repo.GetUserCredentials(ctx, email)
	if err != nil {
		if !errors.Is(err, domainerrors.ErrUserNotFound) {
			log.Printf("failed to get user credentials: %v", err)
		}
		return false, 0
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false, 0
	}

	return true, id
}
