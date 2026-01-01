package services

import (
	"golang-shop-restful/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Services struct {
	db *gorm.DB
}

func NewServices(db *gorm.DB) *Services {
	return &Services{db: db}
}

func (s *Services) GetProducts(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Services) CreateProduct(product *models.Product) error {
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := s.db.Create(product).Error; err != nil {
		return err
	}

	return nil
}

func (s *Services) GetProductById(id int) (models.Product, error) {
	var product models.Product

	if err := s.db.First(&product, id).Error; err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (s *Services) UpdateProduct(id int, updateProduct *models.Product) error {
	var existingProduct models.Product
	if err := s.db.First(&existingProduct, id).Error; err != nil {
		return err
	}

	existingProduct.Title = updateProduct.Title
	existingProduct.Price = updateProduct.Price
	existingProduct.UpdatedAt = time.Now()

	if err := s.db.Save(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}

func (s *Services) DeleteProduct(id int) error {
	var existingProduct models.Product
	if err := s.db.First(&existingProduct, id).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}

// User methods
func (s *Services) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *Services) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Services) AuthenticateUser(username, password string) (models.User, error) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, err
	}

	return user, nil
}
