package main

import (
	"golang-shop-restful/internal/config"
	"golang-shop-restful/internal/database"
	"golang-shop-restful/internal/handlers"
	"golang-shop-restful/internal/middleware"
	"golang-shop-restful/internal/models"
	"golang-shop-restful/internal/services"
	"golang-shop-restful/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoggerInit()
	defer config.Logger.Sync()

	cfg := config.LoadConfig()

	// Initialize JWT with secret from config
	utils.InitJWT(cfg.JWTSecret)

	db, err := database.ConnectDB(cfg)
	if err != nil {
		config.Logger.Fatal(err.Error())
	}
	db.AutoMigrate(&models.Product{}, &models.User{})

	service := services.NewServices(db)
	handler := handlers.NewHandler(service)

	r := gin.Default()
	r.Use(cors.Default())

	auth := r.Group("/api/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)

	api := r.Group("/api/products")
	api.Use(middleware.AuthMiddleware())
	api.GET("/", handler.GetProducts)
	api.POST("/", handler.PostProduct)
	api.GET("/:id", handler.GetProductById)
	api.PUT("/:id", handler.PutProduct)
	api.DELETE("/:id", handler.DeleteProduct)

	config.Logger.Fatal(r.Run())
}
