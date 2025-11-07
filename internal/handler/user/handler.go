package user

import (
	"github.com/WithSoull/UserServer/internal/service"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
	"github.com/WithSoull/platform_common/pkg/tokens"
)

type handler struct {
	desc.UnimplementedUserV1Server
	service       service.UserService
	tokenVerifier tokens.TokenVerifier
}

func NewHandler(service service.UserService, tokenVerifier tokens.TokenVerifier) desc.UserV1Server {
	return &handler{
		service:       service,
		tokenVerifier: tokenVerifier,
	}
}
