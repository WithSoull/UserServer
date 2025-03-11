package config

import (
	"errors"
	"net"
	"os"
)

const (
  grpcHostEnvName = "GRPC_HOST"
  grpcPortEnvName = "GRPC_PORT"
)

type GRPCCongif interface {
  Address() string
}

type grpcConfig struct {
  host string
  port string
}

func NewGRPCConfig() (GRPCCongif, error) {
  host := os.Getenv(grpcHostEnvName)
  if len(host) == 0 {
    return nil, errors.New("grpc host not found")
  }

  port := os.Getenv(grpcPortEnvName)
  if len(port) == 0 {
    return nil, errors.New("grpc port not found")
  }

  return &grpcConfig{
    host: host,
    port: port,
  }, nil
}

func (cfg *grpcConfig) Address() string {
  return net.JoinHostPort(cfg.host, cfg.port)
}
