package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
)

func (r *repo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	builder := sq.Update(usersTableName).
		PlaceholderFormat(sq.Dollar).
		Set(passwordColumn, hashedPassword).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:UpdatePassword",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	log.Print(err)
	return err
}

func (r *repo) LogPassword(ctx context.Context, id int64, ip_address string) error {
	builder := sq.Insert(passwordLogsTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(passwordLogsUserIdColumn, passwordLogsIpAddressColumn).
		Values(id, ip_address)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository:LogPassword",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
