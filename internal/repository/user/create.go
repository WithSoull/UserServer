package user

import (
	"context"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	model "github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/repository/user/conventer"
	"github.com/WithSoull/platform_common/pkg/client/db"
	"github.com/jackc/pgconn"
)

func (r *repo) Create(ctx context.Context, userInfo *model.UserInfo, hashedPassword string, createdAt time.Time) (int64, error) {
	userInfoRepo := conventer.FromModelToRepoUserInfo(userInfo)
	builder := sq.Insert(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, createdAtColumn).
		Values(userInfoRepo.Name, userInfoRepo.Email, hashedPassword, createdAt).
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
