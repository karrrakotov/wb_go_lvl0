package config

import (
	"log"
	"os"

	"github.com/spf13/viper"

	"karrrakotov/wb_go_lvl0/pkg/constants"
)

// Config microservice config
type Config struct {
	AppVersion string
	Nats       Nats
	PostgreSQL PostgreSQL
}

// Nats config
type Nats struct {
	URL       string
	ClusterID string
	ClientID  string
}

// PostgreSQL config
type PostgreSQL struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDBName   string
	PostgresqlSSLMode  string
	PgDriver           string
}

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("internal/config")
	viper.SetConfigName("config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig Parse config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	natsUrl := os.Getenv(constants.NATS_URL)
	if natsUrl != "" {
		c.Nats.URL = natsUrl
	}
	natsClientID := os.Getenv(constants.NATS_CLIENT_ID)
	if natsClientID != "" {
		c.Nats.ClientID = natsClientID
	}
	natsClusterID := os.Getenv(constants.CLUSTER_ID)
	if natsClusterID != "" {
		c.Nats.ClusterID = natsClusterID
	}

	postgresPORT := os.Getenv(constants.POSTGRES_HOST)
	if postgresPORT != "" {
		c.PostgreSQL.PostgresqlHost = postgresPORT
	}

	postgresHost := os.Getenv(constants.POSTGRES_HOST)
	if postgresHost != "" {
		c.PostgreSQL.PostgresqlHost = postgresHost
	}

	postgresqlPort := os.Getenv(constants.POSTGRES_PORT)
	if postgresqlPort != "" {
		c.PostgreSQL.PostgresqlPort = postgresqlPort
	}

	postgresUser := os.Getenv(constants.POSTGRES_USER)
	if postgresUser != "" {
		c.PostgreSQL.PostgresqlUser = postgresUser
	}

	postgresPassword := os.Getenv(constants.POSTGRES_PASSWORD)
	if postgresPassword != "" {
		c.PostgreSQL.PostgresqlPassword = postgresPassword
	}

	postgresDB := os.Getenv(constants.POSTGRES_DB)
	if postgresDB != "" {
		c.PostgreSQL.PostgresqlDBName = postgresDB
	}

	postgresSSL := os.Getenv(constants.POSTGRES_SSL)
	if postgresSSL != "" {
		c.PostgreSQL.PostgresqlSSLMode = postgresSSL
	}

	return &c, nil
}
