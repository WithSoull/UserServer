package user

import (
	model "github.com/WithSoull/AuthService/internal/repository/user/model"
	pb "github.com/WithSoull/AuthService/pkg/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func FromModelToProtoUser(model model.User) *pb.User {
	return &pb.User{
		Id: model.Id,
		Name: model.Name,
		Email: model.Email,
		Role:      pb.Role(model.Role),
		CreatedAt: timestamppb.New(model.CreatedAt),
		UpdatedAt: timestamppb.New(model.UpdatedAt),
	}
}

func FromProtoToModelUser(proto *pb.User) model.User {
	return model.User{
		Id: proto.GetId(),
		Name:  proto.GetName(),
		Email: proto.GetEmail(),
		Role:  model.Role(proto.GetRole()),
		CreatedAt: proto.GetCreatedAt().AsTime(),
		UpdatedAt: proto.GetUpdatedAt().AsTime(),
	}
}


func FromRoleToString(role model.Role) string {
	switch role {
	case model.ROLE_ADMIN:
		return "ADMIN"
	default:
		return "USER"
	}
}

func FromStringToRole(s string) model.Role {
	switch s {
	case "ADMIN":
		return model.ROLE_ADMIN
	default:
		return model.ROLE_USER
	}
}
