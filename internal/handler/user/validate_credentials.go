package user

import (
	"context"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
)

func (h *handler) ValidateCredentials(ctx context.Context, req *desc.ValidateCredentialsRequest) (*desc.ValidateCredentialsResponse, error) {
	valid, id := h.service.ValidateCredentials(ctx, req.GetEmail(), req.GetPassword())
	return &desc.ValidateCredentialsResponse{
		Valid:  valid,
		UserId: id,
	}, nil
}
