package user

import (
	"context"

)

func (s *service) UpdatePassword(ctx context.Context, id int64, password, passwordConfirm string) (error) {
	hashedPassword, err := s.hashPassword(password, passwordConfirm)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, id, hashedPassword)	
}

