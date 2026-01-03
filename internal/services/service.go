package services

import (
	"gorm.io/gorm"
)

type ServicesInterface interface {
	ProductService
	UserService
	CartService
}

type Services struct {
	db *gorm.DB
}

func NewServices(db *gorm.DB) *Services {
	return &Services{db: db}
}
