package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Minio    MinioConfig    `mapstructure:"minio"`
	MQConfig MQConfig       `mapstructure:"mq_config"`
	Backend  BackendConfig  `mapstructure:"backend"`
}

type ServerConfig struct {
	Mode        string `mapstructure:"mode"`
	GrpcPort    int    `mapstructure:"grpc_port"`
	GrpcWebPort int    `mapstructure:"grpc_web_port"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type MinioConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

type ClusterRole string

const (
	ClusterRoleProducer ClusterRole = "producer"
	ClusterRoleConsumer ClusterRole = "consumer"
	ClusterRoleHybrid   ClusterRole = "hybrid" // producer and consumer
)

type BackendConfig struct {
	Jwt            JWTConfig   `mapstructure:"jwt"`
	SystemName     string      `mapstructure:"system_name"`
	ClusterRole    ClusterRole `mapstructure:"cluster_role"`
	InternalRPCKey string      `mapstructure:"internal_rpc_key"`
	// Log           LogConfig      `mapstructure:"log"`
	// document      DocumentConfig `mapstructure:"document"`
}

type MQConfig struct {
	Type     string              `mapstructure:"type"`
	Redis    RedisQueueConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQQueueConfig `mapstructure:"rabbitmq"`
}

type RedisQueueConfig struct {
}

type RabbitMQQueueConfig struct {
	URL string `mapstructure:"url"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int64  `mapstructure:"expire"`
}

var GlobalConfig *Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	return nil
}
