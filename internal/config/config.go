package config

import (
	"os"

	"github.com/WithSoull/UserServer/internal/config/env"
	"github.com/joho/godotenv"
)

// appConfig holds the global application configuration instance.
var appConfig *config

// config represents the complete application configuration.
type config struct {
	Logger LoggerConfig
	GRPC   GRPCConfig
	HTTP   HTTPConfig
	PG     PGConfig
}

// Load reads environment variables from .env file(s) and initializes the application configuration.
// The function ignores file-not-found errors, allowing configuration to be loaded
// from system environment variables when .env file is absent.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	pgCfg, err := env.NewPGConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger: loggerCfg,
		GRPC:   grpcCfg,
		HTTP:   httpCfg,
		PG:     pgCfg,
	}

	return nil
}

// AppConfig returns the global application configuration instance.
// Panics if called before Load().
func AppConfig() *config {
	return appConfig
}
