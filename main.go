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

	api := r.Group("/api/products")
	api.Use(middleware.AuthMiddleware())
	api.GET("/", handler.GetProducts)
	api.POST("/", handler.PostProduct)
	api.GET("/:id", handler.GetProductById)
	api.PUT("/:id", handler.PutProduct)

	cart := r.Group("/api/cart")
	cart.Use(middleware.AuthMiddleware())
	cart.GET("/", handler.GetCart)
	cart.POST("/add/:id", handler.AddToCart)

	utils.Logger.Fatal(r.Run(config.GetServerAddress(cfg)))
}
