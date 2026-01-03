package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler interface {
	GetCart(c *gin.Context)
	AddToCart(c *gin.Context)
	ClearCart(c *gin.Context)
}

func (h *Handler) GetCart(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	cart, err := h.service.GetCart(int(userId.(float64)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get cart",
		})
	}

	c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddToCart(c *gin.Context) {
	productIdStr := c.Param("id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := h.service.AddToCart(int(userId.(float64)), productId); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Product not found or insufficient stock",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to add product to cart",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product added to cart successfully",
	})
}

func (h *Handler) ClearCart(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := h.service.ClearCart(int(userId.(float64))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear cart",
		})
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
