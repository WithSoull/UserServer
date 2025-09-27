package user

import (
	"github.com/WithSoull/UserServer/internal/service"
	desc "github.com/WithSoull/UserServer/pkg/user/v1"
)

type handler struct {
	desc.UnimplementedUserV1Server
	service service.UserService
}

func NewHandler(service service.UserService) desc.UserV1Server {
	return &handler{
		service: service,
	}
}
