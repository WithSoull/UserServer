package user

import (
	"context"
	"log"

	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain_errors"
	"github.com/WithSoull/AuthService/internal/model"
)

func (s *service) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		isLogNeeded, grpcErr := domainerrors.ToGRPCStatus(err)
		if isLogNeeded {
			log.Printf("[Service Layer] failed to get user: %v", err)
		}
		return nil, grpcErr
	}

	return user, nil
}
