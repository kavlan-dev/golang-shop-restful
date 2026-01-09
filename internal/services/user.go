package services

import (
	"golang-shop-restful/internal/models"

	"golang.org/x/crypto/bcrypt"
)

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

	return user, nil
}

func (s *Services) CreateAdminIfNotExists(adminUsername, adminEmail, adminPassword string) error {
	admin, _ := s.getUserByUsername(adminUsername)
	if admin.ID != 0 {
		return nil
	}

	adminUser := &models.User{
		Username: adminUsername,
		Password: adminPassword,
		Email:    adminEmail,
		Role:     "admin",
	}

	if err := s.CreateUser(adminUser); err != nil {
		return err
	}
	if err := s.CreateCart(adminUser); err != nil {
		return err
	}

	return nil
}

func (s *Services) PromoteUserToAdmin(userID int) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == "admin" {
		return nil
	}

	user.Role = "admin"
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func (s *Services) DowngradeUserToCustomer(userID int) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if user.Role == "customer" {
		return nil
	}

	user.Role = "customer"
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
