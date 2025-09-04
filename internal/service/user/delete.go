package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain_errors"
)

func (s *service) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to delete user: %v", err)
		}
		return grpcErr
	}

	return nil
}
