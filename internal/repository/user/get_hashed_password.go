package user

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/UserServer/internal/client/db"
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
	"github.com/jackc/pgx/v4"
)

func (r *repo) GetUserCredentials(ctx context.Context, email string) (int64, string, error) {
	builder := sq.Select(idColumn, passwordColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTableName).
		Where(sq.Eq{emailColumn: email}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, "", err
	}

	q := db.Query{
		Name:     "user_repository:Get",
		QueryRaw: query,
	}

	var (
		hashedPassword string
		id             int64
	)
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", domainerrors.ErrUserNotFound
		}
		return 0, "", err
	}

	return id, hashedPassword, nil
}
