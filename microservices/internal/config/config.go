package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	GrpsServer GrpcServerConfig
	Postgres   PostgresConfig
	Redis      RedisConfig
}

type GrpcServerConfig struct {
	Port         string `env:"PORT" env-default:"8080"`
	JwtSecretKey string `env:"JWT" env-default:"secretkey"`
}

type PostgresConfig struct {
	PostgresqlHost     string `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresqlPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresqlUser     string `env:"POSTGRES_USER" env-default:"postgres"`
	PostgresqlDbname   string `env:"POSTGRES_DB" env-default:"acc_db"`
	PostgresqlPassword string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	PgDriver           string `env:"PGDRIVER" env-default:"pgx"`
}

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
