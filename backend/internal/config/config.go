package config

import (
	"strings"

	"github.com/XingfenD/yoresee_doc/internal/status"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Minio    MinioConfig    `mapstructure:"minio"`
	MQConfig MQConfig       `mapstructure:"mq_config"`
	Consul   ConsulConfig   `mapstructure:"consul"`
	Backend  BackendConfig  `mapstructure:"backend"`
}

type ServerConfig struct {
	GrpcPort    int `mapstructure:"grpc_port"`
	GrpcWebPort int `mapstructure:"grpc_web_port"`
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

type ConsulConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Address    string `mapstructure:"address"`
	Scheme     string `mapstructure:"scheme"`
	Token      string `mapstructure:"token"`
	Datacenter string `mapstructure:"datacenter"`
	Prefix     string `mapstructure:"prefix"`
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
	Jwt            JWTConfig `mapstructure:"jwt"`
	Log            LogConfig `mapstructure:"log"`
	SystemName     string    `mapstructure:"system_name"`
	InternalRPCKey string    `mapstructure:"internal_rpc_key"`
}

type LogConfig struct {
	Level        string `mapstructure:"level"`
	GormLogLevel string `mapstructure:"gorm_log_level"`
}

type MQConfig struct {
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
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "read config failed")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "unmarshal config failed")
	}

	levelText := strings.TrimSpace(GlobalConfig.Backend.Log.Level)
	if levelText == "" {
		levelText = "info"
	}
	level, err := logrus.ParseLevel(strings.ToLower(levelText))
	if err != nil {
		return status.GenErrWithCustomMsg(status.StatusServiceInternalError, "invalid backend.log.level")
	}
	logrus.SetLevel(level)

	return nil
}
