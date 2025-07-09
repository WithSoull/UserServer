package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/WithSoull/AuthService/internal/config"
	"github.com/WithSoull/AuthService/internal/config/env"
	huser "github.com/WithSoull/AuthService/internal/handler/user"
	ruser "github.com/WithSoull/AuthService/internal/repository/user"
	suser "github.com/WithSoull/AuthService/internal/service/user"
	"github.com/WithSoull/AuthService/internal/service"

	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
  flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
  desc.UnimplementedUserV1Server
	userRepository service.UserService
}

func NewServer(service service.UserService) *server {
	return &server{
		userRepository: service,
	}
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
	log.Printf("pgConfig.DSN() = %s", pgConfig.DSN())
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

	userRepo := ruser.NewRepository(dbPool)
	userService := suser.NewService(userRepo)
	userHandler := huser.NewHandler(userService)

	s := grpc.NewServer()
	desc.RegisterUserV1Server(s, userHandler)

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
