package env

import "github.com/IBM/sarama"

type saramaConfig struct {
}

func NewSaramaConfig() (*saramaConfig, error) {
	return &saramaConfig{}, nil
}

func (cfg *saramaConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
