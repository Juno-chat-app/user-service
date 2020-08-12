//
// This file contains configuration of the service
package config

import "time"

type (
	Configuration struct {
		AuthConfig struct {
			AccessTTL  time.Duration `yaml:"auth.accessTTL"`
			RefreshTTL time.Duration `yaml:"auth.refreshTTL"`
		} `yaml:"user_service.auth"`

		GRPCConfig struct {
			Host string `yaml:"host"`
			Port int32  `yaml:"port"`
		} `yaml:"user_service.grpc"`

		CQRSConfig struct {
			PersistConfig PersistConfig `yaml:"persist"`

			CacheConfig struct {
				IsEnable bool   `yaml:"redis.isEnable"`
				Host     string `yaml:"redis.host"`
				Port     int32  `yaml:"redis.port"`
				Db       int    `yaml:"redis.db"`
				Password string `yaml:"redis.password"`
				Retry    int32  `yaml:"redis.retry"`
			} `yaml:"cache"`
		} `yaml:"user_service.cqrs"`
	}

	PersistConfig struct {
		Host              string `yaml:"mongo.host"`
		Port              int32  `yaml:"mongo.port"`
		ConnectionUri     string `yaml:"mongo.uri"`
		UserName          string `yaml:"mongo.userName"`
		Password          string `yaml:"mongo.password"`
		UserDatabase      string `yaml:"mongo.userDatabase"`
		UserCollection    string `yaml:"mongo.userCollection"`
		ConnectionTimeout int64  `yaml:"mongo.connectionTimeout"`
		// todo :: add max pull size and read concern and other options
	}
)
