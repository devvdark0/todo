package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env         string         `env:"ENV" env-default:"local"`
	Port        string         `env:"PORT" env-default:"80"`
	Timeout     time.Duration  `env:"TIMEOUT"`
	IdleTimeout time.Duration  `env:"IDLE_TIMEOUT"`
	DbConfig    DatabaseConfig `env-prefix:"DB_"`
	JWTConfig   JWTConfig
}

type DatabaseConfig struct {
	DBUser     string `env:"USER" env-required:"true"`
	DBPassword string `env:"PASSWORD" env-required:"true"`
	DBHost     string `env:"HOST" env-default:"localhost"`
	DBPort     string `env:"PORT" env-default:"3306"`
	DBName     string `env:"NAME" env-default:"todo-db"`
}

type JWTConfig struct {
	Secret   string        `env:"SECRET_KEY"`
	TokenTTL time.Duration `env:"TOKEN_TTL"`
}

func MustLoad() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
