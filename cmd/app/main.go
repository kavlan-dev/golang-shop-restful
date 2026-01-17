package main

import (
	"go-shop-restful/internal/app"
	"go-shop-restful/internal/config"
	"go-shop-restful/internal/handlers"
	"go-shop-restful/internal/services"
	"go-shop-restful/internal/storage/postgres"
	"go-shop-restful/internal/utils"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := utils.InitLogger(cfg.Environment)
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	utils.InitJWT(cfg.JWTSecret)

	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	storage := postgres.NewStorage(db)
	service := services.NewServices(storage)
	handler := handlers.NewHandler(service, logger)

	if err := service.CreateAdminIfNotExists(cfg.AdminName, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		logger.Errorf("Ошибка создания аккаунта администратора: %v", err)
	}

	log.Fatal(app.Router(cfg, handler))
}
