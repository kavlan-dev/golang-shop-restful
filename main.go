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
		log.Fatal(err.Error())
	}
	defer utils.Logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	utils.InitJWT(cfg.JWTSecret)

	db, err := database.ConnectDB(cfg)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	if err := db.AutoMigrate(&models.Product{}, &models.User{}, &models.Cart{}, &models.CartItem{}); err != nil {
		utils.Logger.Fatal(err.Error())
	}

	service := services.NewServices(db)
	handler := handlers.NewHandler(service)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware(cfg.AllowOrigins))

	auth := r.Group("/api/auth")
	auth.POST("/register", handler.Register)
	auth.POST("/login", handler.Login)

	product := r.Group("/api/products")
	product.Use(middleware.AuthMiddleware())
	product.GET("/", handler.GetProducts)
	product.POST("/", handler.PostProduct)
	product.GET("/:id", handler.GetProductById)
	product.PUT("/:id", handler.PutProduct)
	product.DELETE("/:id", handler.DeleteProduct)

	cart := r.Group("/api/cart")
	cart.Use(middleware.AuthMiddleware())
	cart.GET("/", handler.GetCart)
	cart.POST("/:id", handler.AddToCart)
	cart.DELETE("/", handler.ClearCart)

	utils.Logger.Fatal(r.Run(config.GetServerAddress(cfg)))
}
