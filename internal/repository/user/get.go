package user

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/WithSoull/AuthService/internal/client/db"
	model "github.com/WithSoull/AuthService/internal/model"
	"github.com/WithSoull/AuthService/internal/repository/user/conventer"
	modelRepo "github.com/WithSoull/AuthService/internal/repository/user/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		log.Printf("failed to get user from db: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}

	return conventer.FromRepoToModelUser(&user), nil
}
