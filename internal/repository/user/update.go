package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *repo) Update(ctx context.Context, id int64, name, email *string) error {
	builder := sq.Update(usersTableName).PlaceholderFormat(sq.Dollar)

	if name != nil {
		builder = builder.Set(nameColumn, *name)
	}
	if email != nil {
		builder = builder.Set(emailColumn, *email)
	}
	if name == nil && email == nil {
		return status.Errorf(codes.InvalidArgument, "no fields to update")
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

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
