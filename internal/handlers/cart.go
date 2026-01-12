package handlers

import (
	"golang-shop-restful/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartService interface {
	GetCart(user_id int) (*models.Cart, error)
	AddToCart(user_id, productID int) error
	ClearCart(user_id int) error
}

func (h *Handler) GetCart(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "не авторизован",
		})
		return
	}

	cart, err := h.service.GetCart(int(userId.(float64)))
	if err != nil {
		h.log.Error("Ошибка при выводе корзины:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось получить корзину",
			"details": err,
		})
		return
	}

	h.log.Debugf("Получена корзина пользователя #%d: %v", userId, cart)
	c.JSON(http.StatusOK, cart)
}

func (h *Handler) AddToCart(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "не авторизован",
		})
		return
	}

	if err := h.service.AddToCart(int(userId.(float64)), productId); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "товар не найден",
			})
		} else {
			h.log.Errorf("Ошибка при добавлении товара #%d в корзину: %v", productId, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "не удалось добавить товар в корзину",
				"details": err,
			})
		}
		return
	}

	h.log.Debugf("Товар #%d добавлен в корзину", productId)
	c.JSON(http.StatusOK, gin.H{
		"message": "товар успешно добавлен в корзину",
	})
}

func (h *Handler) ClearCart(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "не авторизован",
		})
		return
	}

	if err := h.service.ClearCart(int(userId.(float64))); err != nil {
		h.log.Error("Ошибка при отчистке корзины:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось отчистить корзину",
			"details": err,
		})
		return
	}

	h.log.Debug("Корзина отчищена")
	c.JSON(http.StatusNoContent, gin.H{})
}
