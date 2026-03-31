package storage

type PostgresOptions struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	GormLogLevel string
}

type RedisOptions struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type MinioOptions struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type ElasticsearchOptions struct {
	Enabled   bool
	Addresses []string
	Username  string
	Password  string
	Timeout   int
}

type ConsulOptions struct {
	Enabled    bool
	Address    string
	Scheme     string
	Token      string
	Datacenter string
	Prefix     string
}
