package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
)

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
