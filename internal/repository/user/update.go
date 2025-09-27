package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/UserServer/internal/client/db"
	domainerrors "github.com/WithSoull/UserServer/internal/errors/domain_errors"
)

func (r *repo) Update(ctx context.Context, id int64, name, email *string) error {
	builder := sq.Update(usersTableName).PlaceholderFormat(squirrel.Dollar)

	if name != nil {
		builder = builder.Set(nameColumn, *name)
	}
	if email != nil {
		builder = builder.Set(emailColumn, *email)
	}
	if name == nil && email == nil {
		return fmt.Errorf("%w: %s", domainerrors.ErrInvalidInput, "no fields to update")
	}

	builder = builder.Set(updatedAtColumn, time.Now())
	builder = builder.Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:Update",
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

	return err
}
