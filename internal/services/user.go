package services

import (
	"golang-shop-restful/internal/models"
	"golang-shop-restful/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user *models.User) error
	AuthenticateUser(username, password string) (models.User, error)
}

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

func (s *Services) getUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *Services) AuthenticateUser(username, password string) (models.User, error) {
	user, err := s.getUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return models.User{}, err
	}

	utils.Logger.Debug(user)
	return user, nil
}
