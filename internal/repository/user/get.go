package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/platform_common/pkg/client/db"
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	model "github.com/WithSoull/UserServer/internal/model"
	"github.com/WithSoull/UserServer/internal/repository/user/conventer"
	modelRepo "github.com/WithSoull/UserServer/internal/repository/user/model"
	"github.com/jackc/pgx/v4"
)

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository:Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}

	return conventer.FromRepoToModelUser(&user), nil
}
