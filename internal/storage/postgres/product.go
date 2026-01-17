package postgres

import "go-shop-restful/internal/models"

func (s *Storage) GetProducts(limit, offset int) (*[]models.Product, error) {
	var products []models.Product
	err := s.db.Limit(limit).Offset(offset).Find(&products).Error
	return &products, err
}

func (s *Storage) CreateProduct(product *models.Product) error {
	return s.db.Create(&product).Error
}

func (s *Storage) FindProductById(id int) (*models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	return &product, err
}

func (s *Storage) UpdateProduct(id int, updateProduct *models.Product) error {
	product, err := s.FindProductById(id)
	if err != nil {
		return err
	}
	return s.db.Model(&product).Updates(&updateProduct).Error
}

func (s *Storage) DeleteProduct(product *models.Product) error {
	return s.db.Delete(&product).Error
}
