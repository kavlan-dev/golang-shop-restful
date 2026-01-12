package handlers

import (
	"golang-shop-restful/internal/models"
	"golang-shop-restful/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateCart(user *models.User) error
	CreateUser(user *models.User) error
	AuthenticateUser(username, password string) (*models.User, error)
	PromoteUserToAdmin(userID int) error
	DowngradeUserToCustomer(userID int) error
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Ошибка в теле запроса регистрации:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "не верное тело запроса",
			"details": err,
		})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := h.service.CreateUser(user); err != nil {
		h.log.Error("Ошибка при создании пользователя:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось создать пользователя",
			"details": err,
		})
		return
	}

	if err := h.service.CreateCart(user); err != nil {
		h.log.Error("Ошибка создания корзины пользователя #%d: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось создать корзину для пользователя",
			"details": err,
		})
		return
	}

	h.log.Debugf("Успешное создание пользователя:", user)
	c.JSON(http.StatusCreated, gin.H{
		"message": "пользователь успешно создан",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Ошибка в теле запроса логина:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "не верное тело запроса",
			"details": err,
		})
		return
	}

	user, err := h.service.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		h.log.Error("Ошибка авторизации:", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "не удалось авторизоваться",
			"details": err,
		})
		return
	}
	h.log.Debug(user)

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		h.log.Error("Ошибка генерации токена:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось сгенерировать JWT токен",
			"details": err,
		})
		return
	}

	h.log.Debugf("Пользователь #%d успешно вошел", user.ID)
	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
	})
}

func (h *Handler) PromoteToAdmin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "не верно введен id",
		})
		return
	}

	if err := h.service.PromoteUserToAdmin(userID); err != nil {
		h.log.Errorf("Ошибка повышения пользователя #%d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "не удалось назначить пользователя администратором",
			"details": err,
		})
		return
	}

	h.log.Debugf("Пользователь #%d повышен", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь повышен до администратора",
	})
}

func (h *Handler) DowngradeToCustomer(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Не верно введен id",
		})
		return
	}

	if err := h.service.DowngradeUserToCustomer(userID); err != nil {
		h.log.Errorf("Ошибка понижения пользователя #%d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Не удалось понизить пользователя",
			"details": err,
		})
		return
	}

	h.log.Debugf("Пользователь #%d понижен", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "пользователь понижен",
	})
}
