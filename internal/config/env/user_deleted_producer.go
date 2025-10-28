package env

import (
	"github.com/caarlos0/env/v11"
)

type userDeletedProducerEnvConfig struct {
	TopicName string `env:"USER_DELETED_TOPIC_NAME,required"`
}

type userDeletedProducerConfig struct {
	raw userDeletedProducerEnvConfig
}

func NewUserDeletedProducerConfig() (*userDeletedProducerConfig, error) {
	var raw userDeletedProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &userDeletedProducerConfig{raw: raw}, nil
}

func (cfg *userDeletedProducerConfig) Topic() string {
	return cfg.raw.TopicName
}
