package services

import "go-shop-restful/internal/models"

type ProductStorage interface {
	GetProducts(limit, offset int) (*[]models.Product, error)
	CreateProduct(product *models.Product) error
	FindProductById(id int) (*models.Product, error)
	UpdateProduct(id int, updateProduct *models.Product) error
	DeleteProduct(product *models.Product) error
}

func (s *Services) GetProducts(limit, offset int) (*[]models.Product, error) {
	return s.storage.GetProducts(limit, offset)
}

func (s *Services) CreateProduct(product *models.Product) error {
	return s.storage.CreateProduct(product)
}

func (s *Services) GetProductById(id int) (*models.Product, error) {
	return s.storage.FindProductById(id)
}

func (s *Services) UpdateProduct(id int, updateProduct *models.Product) error {
	_, err := s.GetProductById(id)
	if err != nil {
		return err
	}
	return s.storage.UpdateProduct(id, updateProduct)
}

func (s *Services) DeleteProduct(id int) error {
	product, err := s.GetProductById(id)
	if err != nil {
		return err
	}

	return s.storage.DeleteProduct(product)
}
