package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/queries"
	// "github.com/WithSoull/AuthService/internal/utils"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	// "google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	dbDSN = "host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable"
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
  return &desc.GetResponse{
    Id: userId,
    Name: gofakeit.Name(),
    Email: gofakeit.Email(),
    Role: 0,

    CreatedAt: timestamppb.New(gofakeit.Date()),
    UpdatedAt: timestamppb.New(gofakeit.Date()),
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
	// if req.GetName() == "" {
	// 		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	// }
	// if req.GetEmail() == "" {
	// 		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	// }
	// if !utils.IsValidEmail(req.GetEmail()) {
	// 		return nil, status.Errorf(codes.InvalidArgument, "invalid email format")
	// }
	// if req.GetPassword() == "" {
	// 		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	// }
	// if req.GetPassword() != req.GetPasswordConfirm() {
	// 		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	// }
	// if req.GetRole() < 0 { // Assuming role is an enum with positive values
	// 		return nil, status.Errorf(codes.InvalidArgument, "invalid role")
	// }

	// Hash password
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	// if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "failed to hash password")
	// }

	// // Begin transaction
	// tx, err := s.db.Begin(ctx)
	// if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "failed to begin transaction")
	// }
	// defer tx.Rollback(ctx)

	// Insert new user
	// var userID int64
	// err = tx.QueryRow(ctx,
	// 		queries.InsertNewUser,
	// 		req.GetName(),
	// 		req.GetEmail(),
	// 		hashedPassword,
	// 		req.GetRole(),
	// ).Scan(&userID)
	// if err != nil {
	// 		return nil, status.Errorf(codes.Internal, "failed to create user")
	// }

	// // Commit transaction
	// if err := tx.Commit(ctx); err != nil {
	// 		return nil, status.Errorf(codes.Internal, "failed to commit transaction")
	// }

	return &desc.CreateResponse{
			Id: 1,
	}, nil

	// return &desc.CreateResponse{
	// 		Id: userID,
	// }, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
  log.Printf("User updating {Id: %d, Name: %s, Email: %s}",
    req.GetId(),
    req.GetName().GetValue(),
    req.GetEmail().GetValue(),
  )

  // TODO: Updating user

  return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
  log.Printf("User deleting {Id: %d}",
    req.GetId(),
  )

  // TODO: Updating user

  return &emptypb.Empty{}, nil
}

func main() {
  flag.Parse()
	ctx := context.Background()	

	
	// Load config
  if err := config.Load(configPath); err != nil {
    log.Fatalf("failed load config: %s", err)
  }

  grpcConfig, err := config.NewGRPCConfig()
  if err != nil {
    log.Fatalf("failed load grpc config: %s", err)
  }

	// Create connection pool
	dbPool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to create connection pool: %s", err)
	}
	defer dbPool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, NewServer(dbPool))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
