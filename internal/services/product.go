package services

import (
	"golang-shop-restful/internal/models"
)

type ProductService interface {
	GetProducts(limit, offset int) ([]models.Product, error)
	CreateProduct(product *models.ProductCreateRequest) error
	GetProductById(id int) (models.Product, error)
	UpdateProduct(id int, updateProduct *models.ProductUpdateRequest) error
	DeleteProduct(id int) error
}

func (s *Services) GetProducts(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	if err := s.db.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *Services) CreateProduct(product *models.ProductCreateRequest) error {
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

func (s *Services) UpdateProduct(id int, updateProduct *models.ProductUpdateRequest) error {
	existingProduct, err := s.GetProductById(id)
	if err != nil {
		return err
	}

	if err := s.db.Model(&existingProduct).Updates(updateProduct).Error; err != nil {
		return err
	}

	return nil
}

func (s *Services) DeleteProduct(id int) error {
	product, err := s.GetProductById(id)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}
