package postgres

import "go-shop-restful/internal/models"

func (s *Storage) CreateUser(user *models.User) error {
	return s.db.Create(&user).Error
}

func (s *Storage) FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Cart").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (s *Storage) FindUserById(userId int) (*models.User, error) {
	var user models.User
	err := s.db.Preload("Cart").First(&user, userId).Error
	return &user, err
}

func (s *Storage) UpdateUser(userId int, updateUser *models.User) error {
	user, err := s.FindUserById(userId)
	if err != nil {
		return err
	}
	return s.db.Model(&user).Updates(updateUser).Error
}
