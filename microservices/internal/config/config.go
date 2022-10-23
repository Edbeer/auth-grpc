package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GrpsServer GrpcServerConfig
	RestServer RestServerConfig
	Postgres   PostgresConfig
	Redis      RedisConfig
	Session    SessionConfig
	Jaeger     JaegerConfig
}

// gRPC server config
type GrpcServerConfig struct {
	Port              string `env:"GRPC_PORT" env-default:"8080"`
	JwtSecretKey      string `env:"JWT" env-default:"secretkey"`
	MaxConnectionIdle int    `env:"MAX_CONNECTION_IDLE" env-default:"10"`
	Timeout           int    `env:"TIMEOUT" env-default:"15"`
	MaxConnectionAge  int    `env:"MAX_CONNECTION_AGE" env-default:"10"`
	Time              int    `env:"TIME" env-default:"120"`
}

type RestServerConfig struct {
	Port         string `env:"REST_PORT" env-default:":9090"`
	JwtSecretKey string `env:"JWT" env-default:"secretkey"`
	ReadTimeout  int    `env:"READ_TIMEOUT" env-default:"10"`
	WriteTimeout int    `env:"WRITE_TIMEOUT" env-default:"10"`
	IdleTimeout  int    `env:"IDLE_TIMEOUT" env-default:"15"`
	TLS          bool   `env:"TLS" env-default:"false"`
}

// Postgres config
type PostgresConfig struct {
	PostgresqlHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresqlPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresqlUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresqlDbname   string `env:"POSTGRES_DB" env-default:"acc_db"`
	PostgresqlPassword string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	PgDriver           string `env:"PGDRIVER" env-default:"pgx"`
}

// Redis config
type RedisConfig struct {
	RedisAddr      string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	RedisPassword  string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB        string `env:"REDIS_DB" env-default:"0"`
	RedisDefaultdb string `env:"REDIS_DDB" env-default:"0"`
	MinIdleConns   int    `env:"MIN_IDLE_CONNS" env-default:"200"`
	PoolSize       int    `env:"POOL_SIZE" env-default:"12000"`
	PoolTimeout    int    `env:"POOL_TIMEOUT" env-default:"240"`
	Password       string `env:"PASSWORD" env-default:""`
	DB             int    `env:"DB" env-default:"0"`
}

// Session config
type SessionConfig struct {
	ExpireAt int `env:"EXPIRE" env-default:"86400"`
}

// Jaeger config
type JaegerConfig struct {
	Host        string `env:"JAEGER_HOST" env-default:"localhost:6831"`
	ServiceName string `env:"SERVICE_NAME" env-default:"ACC_GRPC"`
	LogSpans    bool   `env:"LOG_SPANS" env-default:"false"`
}

var (
	config *Config
	once   sync.Once
)

// Get the config file
func GetConfig() *Config {
	once.Do(func() {
		log.Println("read application configuration")
		config = &Config{}
		if err := cleanenv.ReadEnv(config); err != nil {
			help, _ := cleanenv.GetDescription(config, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return config
}
