package env

import (
	"github.com/caarlos0/env/v11"
)

type userCreatedProducerEnvConfig struct {
	TopicName string `env:"USER_CREATED_TOPIC_NAME,required"`
}

type userCreatedProducerConfig struct {
	raw userCreatedProducerEnvConfig
}

func NewUserCreatedProducerConfig() (*userCreatedProducerConfig, error) {
	var raw userCreatedProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &userCreatedProducerConfig{raw: raw}, nil
}

func (cfg *userCreatedProducerConfig) Topic() string {
	return cfg.raw.TopicName
}
