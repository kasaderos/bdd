package config

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	Templates map[string][]string
}

type DatabaseConfig struct {
	Name     string `envconfig:"DB_NAME"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Host     string `envconfig:"DB_HOST"`
	Port     string `envconfig:"DB_PORT"`
}

func (d DatabaseConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
	)
}

type Config struct {
	LogLevel           string `envconfig:"LOG_LEVEL"`
	HTTPListenPort     string `envconfig:"HTTP_LISTEN_PORT"`
	AnalyticsSystemURL string `envconfig:"ANALYTICS_SYSTEM_URL"`

	DatabaseConfig DatabaseConfig
	ServiceConfig  ServiceConfig
}

func Load() *Config {
	var c Config

	if err := envconfig.Process("", &c); err != nil {
		log.Fatal("config ", err)
	}

	c.ServiceConfig = loadServiceConfig()
	c.DatabaseConfig = loadDatabaseConfig()

	return &c
}

func loadServiceConfig() ServiceConfig {
	var c ServiceConfig
	if err := envconfig.Process("SERVICE", &c); err != nil {
		log.Fatal("service config ", err)
	}

	// templates mapping
	c.Templates = map[string][]string{
		"main":        {"base"},
		"admin":       {"admin"},
		"asset_table": {"asset_table"},
		"chart":       {"chart"},
	}

	return c
}

func loadDatabaseConfig() DatabaseConfig {
	var conf DatabaseConfig
	if err := envconfig.Process("DB", &conf); err != nil {
		log.Fatal("db config ", err)
	}

	return conf
}
