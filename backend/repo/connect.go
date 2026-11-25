package repo

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/zakirab-git/stockflow/config"
)

func Connect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sqlx.Open("pgx", connectStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
