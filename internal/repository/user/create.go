package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
	domainerrors "github.com/WithSoull/AuthService/internal/errors/domain_errors"
	model "github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/repository/user/conventer"
	"github.com/jackc/pgconn"
)

func (r *repo) Create(ctx context.Context, userInfo *model.UserInfo, hashedPassword string) (int64, error) {
	userInfoRepo := conventer.FromModelToRepoUserInfo(userInfo)
	builder := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn).
		Values(userInfoRepo.Name, userInfoRepo.Email, hashedPassword).
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, domainerrors.ErrEmailAlreadyExists
		}
		return 0, err
	}

	return userID, nil
}
