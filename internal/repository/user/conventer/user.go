package conventer

import (
	"github.com/WithSoull/AuthService/internal/model"
	modelRepo "github.com/WithSoull/AuthService/internal/repository/user/model"
)

func FromRepoToModelUserInfo(userInfo *modelRepo.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name: userInfo.Name,
		Email: userInfo.Email,
		Role: FromStringToRole(userInfo.Role),
	}	
}

func FromModelToRepoUserInfo(userInfo *model.UserInfo) *modelRepo.UserInfo {
	return &modelRepo.UserInfo{
		Name: userInfo.Name,
		Email: userInfo.Email,
		Role: FromRoleToString(userInfo.Role),
	}	
}

func FromRepoToModelUser(user *modelRepo.User) *model.User {
	return &model.User{
		Id: user.Id,
		UserInfo: *FromRepoToModelUserInfo(&user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func FromModelToRepoUser(user *model.User) *modelRepo.User {
	return &modelRepo.User{
		Id: user.Id,
		UserInfo: *FromModelToRepoUserInfo(&user.UserInfo),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
