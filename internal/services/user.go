package services

import (
	"golang-shop-restful/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	CreateUser(user *models.User) error
	FindUserByUsername(username string) (*models.User, error)
	FindUserById(userId int) (*models.User, error)
	UpdateUser(userId int, updateUser *models.User) error
}

func (s *Services) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return s.storage.CreateUser(user)
}

func (s *Services) getUserByUsername(username string) (*models.User, error) {
	return s.storage.FindUserByUsername(username)
}

func (s *Services) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Services) GetUserById(userId int) (*models.User, error) {
	return s.storage.FindUserById(userId)
}

func (s *Services) UpdateUser(userId int, updateUser *models.User) error {
	return s.storage.UpdateUser(userId, updateUser)
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

func (s *Services) PromoteUserToAdmin(userId int) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Role == "admin" {
		return nil
	}

	user.Role = "admin"
	if err := s.UpdateUser(userId, user); err != nil {
		return err
	}

	return nil
}

func (s *Services) DowngradeUserToCustomer(userId int) error {
	user, err := s.GetUserById(userId)
	if err != nil {
		return err
	}

	if user.Role == "customer" {
		return nil
	}

	user.Role = "customer"
	if err := s.UpdateUser(userId, user); err != nil {
		return err
	}

	return nil
}
