package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/zakirab-git/stockflow/config"
	"github.com/zakirab-git/stockflow/handler"
	"github.com/zakirab-git/stockflow/repo"
	"github.com/zakirab-git/stockflow/router"
	"github.com/zakirab-git/stockflow/server"
)

func main() {
	cfg := config.Load()

	db := initDB(*cfg)
	defer db.Close()

	handlers, err := initHandlers(*cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	server.StartHTTPServer(router.SetupRouter(*cfg, handlers), cfg.App)

}

func initDB(cfg config.Config) *sqlx.DB {
	db, err := repo.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}

func initHandlers(cfg config.Config, db *sqlx.DB) (*handler.Handlers, error) {
	return &handler.Handlers{}, nil
}
