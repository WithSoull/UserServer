package user

import (
	"context"

	conventer "github.com/WithSoull/AuthService/internal/conventer/user"
	"github.com/WithSoull/AuthService/internal/service"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	userID, err := s.service.Create(ctx, conventer.FromProtoToModelUser(req.GetUser()), req.GetPassword(), req.GetPasswordConfirm())
	if err != nil {
		return nil, err
	}

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

func (s *handler) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &desc.GetResponse{
		User: conventer.FromModelToProtoUser(*user), 
	}, nil
}

func (s *handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	var name *string
	if req.GetName() != nil {
		name = &req.GetName().Value
	}

	var email *string
	if req.GetEmail() != nil {
		email = &req.GetEmail().Value
	}

	return &emptypb.Empty{}, s.service.Update(ctx, req.GetId(), name, email)
}


func (s *handler) UpdatePassword(ctx context.Context, req *desc.UpdatePasswordRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.service.UpdatePassword(ctx, req.GetId(), req.GetPassword(), req.GetPasswordConfirm())
}

func (s *handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.service.Delete(ctx, req.GetId())
}
