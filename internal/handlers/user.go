package handlers

import (
	"golang-shop-restful/internal/models"
	"golang-shop-restful/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(user *models.User) error
	AuthenticateUser(username, password string) (models.User, error)
	CreateAdminIfNotExists(adminUsername, adminEmail, adminPassword string) error
	PromoteUserToAdmin(userID int) error
	DowngradeUserToCustomer(userID int) error
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := h.service.CreateUser(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := h.service.CreateCart(user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user successfully created",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := h.service.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
	})
}

func (h *Handler) PromoteToAdmin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.PromoteUserToAdmin(userID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
	})
}

func (h *Handler) DowngradeToCustomer(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := h.service.DowngradeUserToCustomer(userID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User downgraded to customer successfully",
	})
}
