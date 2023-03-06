package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DB         *DB    `envconfig:"db"`
	App        *App   `envconfig:"app"`
	JWTKeyword string `envconfig:"jwt_keyword"`
	Minio      *Minio `envconfig:"minio"`
}

type App struct {
	Port string `envconfig:"port"`
}

type DB struct {
	Driver   string `envconfig:"driver" required:"true"`
	Password string `envconfig:"password"`
	User     string `envconfig:"user"`
	Name     string `envconfig:"name"`
	SSLMode  string `envconfig:"sslmode"`
}

type Minio struct {
	KeyID     string `envconfig:"key_id"`
	SecretKey string `envconfig:"secret_key"`
	Endpoint  string `envconfig:"endpoint"`
	Bucket    string `envconfig:"bucket"`
}

func (c *Config) Process() error {
	return envconfig.Process("example", c)
}
