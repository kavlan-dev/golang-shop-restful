package services

import (
	"golang-shop-restful/internal/models"

	"gorm.io/gorm"
)

type CartService interface {
	CreateCart(user *models.User) error
	GetCart(user_id int) (models.Cart, error)
	AddToCart(user_id, productID int) error
}

func (s *Services) CreateCart(user *models.User) error {
	if user.Cart.UserID == 0 {
		user.Cart = models.Cart{UserID: user.ID}
		if err := s.db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Services) GetCart(user_id int) (models.Cart, error) {
	var cart models.Cart
	if err := s.db.Preload("Items").Where("user_id = ?", user_id).First(&cart).Error; err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

func (s *Services) AddToCart(user_id, productID int) error {
	var user models.User
	if err := s.db.Preload("Cart").First(&user, user_id).Error; err != nil {
		return err
	}

	cart := user.Cart
	if cart.ID == 0 {
		return gorm.ErrRecordNotFound
	}

	var product models.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return err
	}

	if product.Stock <= 0 {
		return gorm.ErrRecordNotFound
	}

	var existingCartItem models.CartItem
	if err := s.db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&existingCartItem).Error; err == nil {
		existingCartItem.Quantity += 1
		existingCartItem.Price = product.Price * float64(existingCartItem.Quantity)
		if err := s.db.Save(&existingCartItem).Error; err != nil {
			return err
		}
	} else {
		newCartItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: uint(productID),
			Quantity:  1,
			Price:     product.Price,
		}
		if err := s.db.Create(&newCartItem).Error; err != nil {
			return err
		}
	}
	product.Stock -= 1
	if err := s.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
