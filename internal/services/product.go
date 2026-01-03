package services

import (
	"golang-shop-restful/internal/models"
)

type ProductService interface {
	GetProducts(limit, offset int) ([]models.Product, error)
	CreateProduct(product *models.Product) error
	GetProductById(id int) (models.Product, error)
	UpdateProduct(id int, updateProduct *models.Product) error
}

func (s *Services) GetProducts(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Services) CreateProduct(product *models.Product) error {
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

	existingProduct = models.Product{
		Title:       updateProduct.Title,
		Description: updateProduct.Description,
		Price:       updateProduct.Price,
		Category:    existingProduct.Category,
		Stock:       existingProduct.Stock,
	}

	if err := s.db.Save(&existingProduct).Error; err != nil {
		return err
	}

	return nil
}
