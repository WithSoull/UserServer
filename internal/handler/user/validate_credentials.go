package user

import (
	"context"
	"strings"

	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/contextx/claimsctx"
	"github.com/WithSoull/platform_common/pkg/sys"
	"github.com/WithSoull/platform_common/pkg/sys/codes"
	"google.golang.org/grpc/metadata"
)

func (h *handler) ValidateCredentials(ctx context.Context, req *desc.ValidateCredentialsRequest) (*desc.ValidateCredentialsResponse, error) {
	valid, id := h.service.ValidateCredentials(ctx, req.GetEmail(), req.GetPassword())
	return &desc.ValidateCredentialsResponse{
		Valid:  valid,
		UserId: id,
	}, nil
}

// verifyToken extract token from grpc metadata
// and then verify token by tokenVerifier
// and then put user's claims in context
func (h *handler) verifyToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, sys.NewCommonError("metadata not provided", codes.Unauthenticated)
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return ctx, sys.NewCommonError("authorization header not provided", codes.Unauthenticated)
	}

	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := h.tokenVerifier.VerifyAccessToken(ctx, token)
	if err != nil {
		return ctx, sys.NewCommonError("invalid token", codes.Unauthenticated)
	}

	ctxWithEmail := claimsctx.InjectUserEmail(ctx, claims.Email)
	ctxWithEmailAndUserID := claimsctx.InjectUserID(ctxWithEmail, claims.UserId)

	return ctxWithEmailAndUserID, nil
}
