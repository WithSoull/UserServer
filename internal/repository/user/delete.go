package user

import (
	"context"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/UserServer/internal/client/db"
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
)

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(usersTableName).
		PlaceholderFormat(squirrel.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:Delete",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domainerrors.ErrUserNotFound
	}

	return nil
}
