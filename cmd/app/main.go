package main

import (
	"golang-shop-restful/internal/config"
	"golang-shop-restful/internal/handlers"
	"golang-shop-restful/internal/middleware"
	"golang-shop-restful/internal/services"
	"golang-shop-restful/internal/storage/postgres"
	"golang-shop-restful/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	logger, err := utils.InitLogger()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(err)
	}

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

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg.AllowOrigins))

	auth := r.Group("/api/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	admin := api.Group("/admin")
	admin.Use(middleware.AdminMiddleware())

	product := api.Group("/products")
	product.GET("/", handler.GetProducts)
	product.GET("/:id", handler.GetProductById)

	cart := api.Group("/cart")
	cart.GET("/", handler.GetCart)
	cart.POST("/:id", handler.AddToCart)
	cart.DELETE("/", handler.ClearCart)

	adminUsers := admin.Group("/users")
	adminUsers.POST("/:id/promote", handler.PromoteToAdmin)
	adminUsers.POST("/:id/downgrade", handler.DowngradeToCustomer)

	adminProduct := admin.Group("/products")
	adminProduct.POST("/", handler.PostProduct)
	adminProduct.PUT("/:id", handler.PutProduct)
	adminProduct.DELETE("/:id", handler.DeleteProduct)

	logger.Fatal(r.Run(config.GetServerAddress(cfg)))
}
