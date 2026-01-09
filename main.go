package main

import (
	"golang-shop-restful/internal/config"
	"golang-shop-restful/internal/database"
	"golang-shop-restful/internal/handlers"
	"golang-shop-restful/internal/middleware"
	"golang-shop-restful/internal/models"
	"golang-shop-restful/internal/services"
	"golang-shop-restful/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := utils.InitLogger(); err != nil {
		log.Fatal(err)
	}
	defer utils.Logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Logger.Fatal(err)
	}

	utils.InitJWT(cfg.JWTSecret)

	db, err := database.ConnectDB(cfg)
	if err != nil {
		utils.Logger.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Product{}, &models.User{}, &models.Cart{}, &models.CartItem{}); err != nil {
		utils.Logger.Fatal(err)
	}

	service := services.NewServices(db)
	handler := handlers.NewHandler(service)

	if err := service.CreateAdminIfNotExists(cfg.AdminName, cfg.AdminEmail, cfg.AdminPassword); err != nil {
		utils.Logger.Errorf("Failed to create initial admin: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg.AllowOrigins))
	r.Use(gin.ErrorLogger())

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

	utils.Logger.Fatal(r.Run(config.GetServerAddress(cfg)))
}
