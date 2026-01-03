package handlers

import (
	"golang-shop-restful/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler interface {
	GetProducts(c *gin.Context)
	PostProduct(c *gin.Context)
	GetProductById(c *gin.Context)
	PutProduct(c *gin.Context)
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
				"error": "Invalid limit parameter",
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
				"error": "Invalid offset parameter",
			})
			return
		}
	}

	products, err := h.service.GetProducts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) PostProduct(c *gin.Context) {
	var newProduct models.Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product data",
		})
		return
	}

	if err := h.service.CreateProduct(&newProduct); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, newProduct)
}

func (h *Handler) GetProductById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := h.service.GetProductById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve product",
			})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) PutProduct(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var updateProduct models.Product
	if err := c.BindJSON(&updateProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product data",
		})
		return
	}

	if err := h.service.UpdateProduct(id, &updateProduct); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update product",
			})
		}
		return
	}

	c.JSON(http.StatusOK, updateProduct)
}
