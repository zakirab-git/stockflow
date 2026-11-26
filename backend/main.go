package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/zakirab-git/stockflow/config"
	"github.com/zakirab-git/stockflow/handler"
	"github.com/zakirab-git/stockflow/repo"
	"github.com/zakirab-git/stockflow/router"
	"github.com/zakirab-git/stockflow/server"
	"github.com/zakirab-git/stockflow/usecase"
	"github.com/zakirab-git/stockflow/utils"
)

func main() {
	cfg := config.Load()

	db := initDB(*cfg)
	defer db.Close()

	handlers := initHandlers(db)

	server.StartHTTPServer(router.SetupRouter(*cfg, handlers), cfg.App)

}

func initDB(cfg config.Config) *sqlx.DB {
	db, err := repo.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}

func initHandlers(db *sqlx.DB) *handler.Handlers {
	txM := utils.NewTxM(db)

	productRepo := repo.NewProductRepo(db)
	productUsecase := usecase.NewProductUsecase(productRepo, txM)
	productHandler := handler.NewProductHandler(productUsecase)

	orderRepo := repo.NewOrderRepo(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, txM)
	orderHandler := handler.NewOrderHandler(orderUsecase)

	return handler.NewHandlers(productHandler, orderHandler)
}
