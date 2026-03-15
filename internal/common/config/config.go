package config

import "time"

// Config represents the main structure of global config
type Config struct {
	App         App         `yaml:"app" validate:"required"`
	Server      Server      `yaml:"server" validate:"required"`
	Persistence Persistence `yaml:"persistence" validate:"required"`
	Cache       Cache       `yaml:"cache" validate:"required"`
}

type App struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required,semver"`
	Env     string `yaml:"env" validate:"required,oneof=dev local prod"`
}

type Server struct {
	Addr         string        `yaml:"addr" validate:"required,hostname_port"`
	TLS          TLSConfig     `yaml:"tls" validate:"required"`
	ReadTimeout  time.Duration `yaml:"read-timeout" validate:"required,min=100ms"`
	WriteTimeout time.Duration `yaml:"write-timeout" validate:"required,min=100ms"`
	IdleTimeout  time.Duration `yaml:"idle-timeout" validate:"required,min=1s"`
}

type TLSConfig struct {
	Enable         bool   `yaml:"enable"`
	ServerCertPath string `yaml:"server-cert-path"`
	ServerKeyPath  string `yaml:"server-key-path"`
}

type Persistence struct {
	Postgres Postgres `yaml:"postgres" validate:"required"`
}

type PostgresAuth struct {
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type PostgresConn struct {
	MaxIdles    int           `yaml:"max-idles" validate:"required,gte=1"`
	MaxOpens    int           `yaml:"max-opens" validate:"required,gte=1"`
	MaxIdleTime time.Duration `yaml:"max-idle-time" validate:"required,min=100ms"`
	MaxLifetime time.Duration `yaml:"max-lifetime" validate:"required,min=100ms"`
}

type Postgres struct {
	Host    string       `yaml:"host" validate:"required,hostname"`
	Port    int          `yaml:"port" validate:"required,number,min=1,max=65535"`
	SSLMode string       `yaml:"sslmode" validate:"required,oneof=disable enable"`
	DBName  string       `yaml:"dbname" validate:"required"`
	Auth    PostgresAuth `yaml:"auth" validate:"required"`
	Conn    PostgresConn `yaml:"conn" validate:"required"`
}

type Cache struct {
	Redis Redis `yaml:"redis" validate:"required"`
}

type RedisAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Redis struct {
	Host string    `yaml:"host" validate:"required,hostname"`
	Port int       `yaml:"port" validate:"required,number,min=1,max=65535"`
	DB   int       `yaml:"db" validate:"gte=0"`
	Auth RedisAuth `yaml:"auth" validate:"required"`
}
