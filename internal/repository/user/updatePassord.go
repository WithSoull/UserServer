package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
)

func (r *repo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Set(passwordColumn, hashedPassword).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:UpdatePassword",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args)
	return err
}
