package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/WithSoull/AuthService/internal/config"
	desc "github.com/WithSoull/AuthService/pkg/user/v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

func init() {
  flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

const (
  grpcPort = 50051
)


type server struct {
  desc.UnimplementedUserV1Server
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

  // TODO: Creating User

  return &desc.CreateResponse{
    Id: gofakeit.Int64(),
  }, nil
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
  
  if err := config.Load(configPath); err != nil {
    log.Fatalf("failed load config: %s", err)
  }

  grpcConfig, err := config.NewGRPCConfig()
  if err != nil {
    log.Fatalf("failed load grpc config: %s", err)
  }

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
