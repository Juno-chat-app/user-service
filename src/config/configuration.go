//
// This file contains configuration of the service
package config

type (
	Configuration struct {
		GRPCConfig struct {
			Host string `yaml:"host"`
			Port int32  `yaml:"port"`
		} `yaml:"user_service.grpc"`

		CQRSConfig struct {
			PersistConfig struct {
				Host string `yaml:"mongo.host"`
				Port int32  `yaml:"mongo.port"`
			} `yaml:"persist"`

			CacheConfig struct {
				IsEnable bool   `yaml:"redis.isEnable"`
				Host     string `yaml:"redis.host"`
				Port     int32  `yaml:"redis.port"`
			} `yaml:"cache"`
		} `yaml:"user_service.cqrs"`
	}
)
