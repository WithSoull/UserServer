package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	conventer "github.com/WithSoull/AuthService/internal/conventer/user"
	"github.com/WithSoull/AuthService/internal/queries"
	"github.com/WithSoull/AuthService/internal/repository"
	"github.com/WithSoull/AuthService/internal/repository/user/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}	
}

func (r *repo) Create(ctx context.Context, user *model.User, hashedPassword string)  (int64, error) {
	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return 0, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert new user
	var userID int64
	err = tx.QueryRow(ctx,
			queries.InsertNewUser,
			user.Name,
			user.Email,
			hashedPassword,
			conventer.FromRoleToString(user.Role),
	).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return 0, status.Errorf(codes.AlreadyExists, "this email already used")
		} else {
			log.Printf("failed to insert user in db: %v", err)
			return 0, status.Errorf(codes.Internal, "failed to create user")
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return 0, status.Errorf(codes.Internal, "failed to commit transaction")
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	var (
		name       string
		email      string
		roleStr    string
		createdAt  time.Time
		updatedAt  time.Time
	)
	err := r.db.QueryRow(ctx, queries.SelectById, id).Scan(
		&name,
		&email,
		&roleStr,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		log.Printf("failed to get user from db: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user")
	}

  return &model.User{
		Id: id,
    Name: name,
    Email: email,
		Role: conventer.FromStringToRole(roleStr),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
  }, nil
}

func (r *repo) Update(ctx context.Context, id int64, name, email *string) (error) {
	// Create query for user updating
	params := make([]interface{}, 0)
	setClauses := make([]string, 0)
	if name != nil {
		name := name
		params = append(params, name)
		clause := fmt.Sprintf("name = %s%d", queries.PlaceHolder, len(params))
		setClauses = append(setClauses, clause)
	}
	if email != nil {
		params = append(params, email)
		clause := fmt.Sprintf("email = %s%d", queries.PlaceHolder, len(params))
		setClauses = append(setClauses, clause)
	}
	if len(params) == 0 {
		return status.Errorf(codes.InvalidArgument, "no fields to update")
	}

	params = append(params, id)
	query := fmt.Sprintf(queries.UpdateById, strings.Join(setClauses, ", "), len(params))

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Update user
	res, err := tx.Exec(ctx, query, params...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		log.Printf("query: %s", query)
		for _, param := range params {
			log.Println(param)
		}
		return status.Errorf(codes.Internal, "failed to update user")
	}
	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "user(%d) not found", id)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return nil
}

func (r *repo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	clause := fmt.Sprintf("password = %s1", queries.PlaceHolder)
	query := fmt.Sprintf(queries.UpdateById, clause, 2)

	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Update password
	res, err := tx.Exec(ctx, query, hashedPassword, id)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return status.Errorf(codes.Internal, "failed to update user")
	}
	if res.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "user(%d) not found", id)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return nil
}

func (r *repo) Delete(ctx context.Context, id int64) (error) {
	// Begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Deleting user
	_, err = tx.Exec(ctx, queries.DeleteById, id)
	if err != nil {
    log.Printf("failed to delete user: %v", err)
		return status.Errorf(codes.Internal, "failed to delete user")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return nil
}

