package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"karrrakotov/wb_go_lvl0/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password='%s' sslmode=%s",
		cfg.PostgreSQL.PostgresqlHost,
		cfg.PostgreSQL.PostgresqlPort,
		cfg.PostgreSQL.PostgresqlUser,
		cfg.PostgreSQL.PostgresqlDBName,
		cfg.PostgreSQL.PostgresqlPassword,
		cfg.PostgreSQL.PostgresqlSSLMode))

	if err != nil {
		log.Fatalln("не удалось подключиться к базе данных")
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalln("не удалось выполнить ping к базе данных")
		return nil, err
	}

	return db, nil
}
