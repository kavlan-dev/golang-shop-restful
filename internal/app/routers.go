package app

import (
	"fmt"
	"go-shop-restful/internal/config"
	"go-shop-restful/internal/handlers"
	"go-shop-restful/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Router(cfg *config.Config, handler *handlers.Handler) error {
	var r *gin.Engine
	switch cfg.Environment {
	case "dev":
		r = gin.Default()
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Logger(), gin.Recovery())
	default:
		return fmt.Errorf("Не известное окружение %s", cfg.Environment)
	}
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

	return r.Run(config.GetServerAddress(cfg))
}
