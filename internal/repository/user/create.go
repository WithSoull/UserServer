package user

import (
	"context"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
	model "github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/repository/user/conventer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *repo) Create(ctx context.Context, userInfo *model.UserInfo, hashedPassword string) (int64, error) {
	userInfoRepo := conventer.FromModelToRepoUserInfo(userInfo)
	builder := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn).
		Values(userInfoRepo.Name, userInfoRepo.Email, hashedPassword, userInfoRepo.Role).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository:Create",
		QueryRaw: query,
	}

	var userID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return 0, status.Errorf(codes.AlreadyExists, "this email already used")
		} else {
			log.Printf("failed to insert user in db: %v", err)
			return 0, status.Errorf(codes.Internal, "failed to create user")
		}
	}

	return userID, nil
}
