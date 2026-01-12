package handlers

import (
	"golang-shop-restful/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductService interface {
	GetProducts(limit, offset int) (*[]models.Product, error)
	CreateProduct(product *models.Product) error
	GetProductById(id int) (*models.Product, error)
	UpdateProduct(id int, updateProduct *models.Product) error
	DeleteProduct(id int) error
}

func (h *Handler) GetProducts(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	var err error
	limit := 100
	offset := 0
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "не верно введен limit",
			})
			return
		}
		if limit > 1000 {
			limit = 1000
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "не верно введен offset",
			})
			return
		}
	}

	products, err := h.service.GetProducts(limit, offset)
	if err != nil {
		h.log.Error("Ошибка при получении всех товаров:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось получить товары",
			"details": err,
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) PostProduct(c *gin.Context) {
	var req models.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Ошибка в теле создания товара:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "не верное тело запроса",
			"details": err,
		})
		return
	}

	newProduct := &models.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
	}

	if err := h.service.CreateProduct(newProduct); err != nil {
		h.log.Errorf("Ошибка создания товара %s: %v", newProduct.Title, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось создать товар",
			"details": err,
		})
		return
	}

	h.log.Debugf("Создан товар #%d с названием %s", newProduct.ID, newProduct.Title)
	c.JSON(http.StatusCreated, newProduct)
}

func (h *Handler) GetProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не верно введен id",
		})
		return
	}

	product, err := h.service.GetProductById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "товар не найден",
			})
		} else {
			h.log.Errorf("Товар #%d не найден: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "не удалось найти товар",
				"details": err,
			})
		}
		return
	}

	h.log.Debugf("Получен товар #%d с названием %s", id, product.Title)
	c.JSON(http.StatusOK, product)
}

func (h *Handler) PutProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не верно введен id",
		})
		return
	}

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Errorf("Ошибка в теле запроса для обновления товара #%d: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "не верное тело запроса",
			"details": err,
		})
		return
	}

	updateProduct := &models.Product{
		Category:    req.Category,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := h.service.UpdateProduct(id, updateProduct); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "товар не найден",
			})
		} else {
			h.log.Errorf("Ошибка при изменении товара #%d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "не удалось обновить товар",
				"details": err,
			})
		}
		return
	}

	h.log.Debugf("Товар #%d с названием %s успешно изменен", id, updateProduct.Title)
	c.JSON(http.StatusOK, gin.H{
		"message": "продукт успешно изменен",
	})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "не верно введен id",
		})
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "товар не найден",
			})
		} else {
			h.log.Error("Ошибка при удалении товара #%d: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "не удалось удалить товар",
				"details": err,
			})
		}
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"message": "продукт успешно удален",
	})
}
