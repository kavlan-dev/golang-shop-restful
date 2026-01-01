package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title string `json:"title" binding:"required,min=3,max=100"`
	Price uint   `json:"price" binding:"required,min=0"`
}
