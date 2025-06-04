package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	"github.com/WithSoull/AuthService/internal/queries"
	"github.com/WithSoull/AuthService/internal/utils"

	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
  flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
  desc.UnimplementedUserV1Server
	db *pgxpool.Pool
}

func NewServer(dbPool *pgxpool.Pool) *server {
	return &server{
		db: dbPool,
	}
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
  userId := req.GetId()
  log.Printf("User getting {Id: %d}",
    req.GetId(),
  )
	
		
	var (
		name       string
		email      string
		role       string
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := s.db.QueryRow(ctx, queries.SelectById, userId).Scan(
		&name,
		&email,
		&role,
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

  return &desc.GetResponse{
    Id: userId,
    Name: name,
    Email: email,
    Role: desc.Role(desc.Role_value[role]),

    CreatedAt: timestamppb.New(createdAt),
    UpdatedAt: timestamppb.New(updatedAt),
  }, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
  log.Printf("User creating {Name: %s, Email: %s, Password: %s, PasswordConfirm: %s, Role: %d}",
    req.GetName(),
    req.GetEmail(), 
    req.GetPassword(), 
    req.GetPasswordConfirm(),
    req.GetRole(),
  )

	// Input validation
	if req.GetName() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.GetEmail() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if !utils.IsValidEmail(req.GetEmail()) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid email format")
	}
	if req.GetPassword() == "" {
			return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}
	if req.GetPassword() != req.GetPasswordConfirm() {
			return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}
	if req.GetRole() < 0 && req.GetRole() > 1 { // Assuming role is an enum with positive values
			return nil, status.Errorf(codes.InvalidArgument, "invalid role")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Insert new user
	var userID int64
	err = tx.QueryRow(ctx,
			queries.InsertNewUser,
			req.GetName(),
			req.GetEmail(),
			hashedPassword,
			req.GetRole(),
	).Scan(&userID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return nil, status.Errorf(codes.AlreadyExists, "this email already used")
		} else {
			log.Printf("failed to insert user in db: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to create user")
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}

	return &desc.CreateResponse{
			Id: userID,
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
  log.Printf("User updating {Id: %d, Name: %s, Email: %s}",
    req.GetId(),
    req.GetName().GetValue(),
    req.GetEmail().GetValue(),
  )
	
	// Create query for user updating
	params := make([]interface{}, 0)
	setClauses := make([]string, 0)
	if req.GetName() != nil {
		name := req.GetName().GetValue()
		params = append(params, name)
		clause := fmt.Sprintf("name = %s%d", queries.PlaceHolder, len(params))
		setClauses = append(setClauses, clause)
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		params = append(params, email)
		clause := fmt.Sprintf("email = %s%d", queries.PlaceHolder, len(params))
		setClauses = append(setClauses, clause)
	}
	if len(params) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no fields to update")
	}

	params = append(params, req.GetId())
	query := fmt.Sprintf(queries.UpdateById, strings.Join(setClauses, ", "), len(params))

	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
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
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}
	if res.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "user(%d) not found", req.GetId())
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
  log.Printf("User deleting {Id: %d}",
    req.GetId(),
  )
	
	// Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
    log.Printf("failed to begin transaction: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Update user
	_, err = tx.Exec(ctx, queries.DeleteById, req.GetId())
	if err != nil {
    log.Printf("failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
    log.Printf("failed to commit transaction: %v", err)
    return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	}

  return &emptypb.Empty{}, nil
}

func main() {
  flag.Parse()
	ctx := context.Background()	

	// Load config
  if err := config.Load(configPath); err != nil {
    log.Printf("configPath=%s", configPath)
    log.Fatalf("failed load config: %s", err)
  }

  grpcConfig, err := env.NewGRPCConfig()
  if err != nil {
    log.Fatalf("failed load grpc config: %s", err)
  }

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed load pg congig: %s", err)
	}

	// Create connection pool
	dbPool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to create connection pool: %s", err)
	}
	defer dbPool.Close()

  if err := dbPool.Ping(ctx); err != nil {
    log.Fatalf("failed to connect to database: %v", err)
  }

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterUserV1Server(s, NewServer(dbPool))

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
