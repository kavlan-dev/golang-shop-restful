package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"size:50;unique;not null" binding:"required,min=3,max=50"`
	Password string `json:"password" gorm:"size:100;not null" binding:"required,min=6,max=100"`
	Email    string `json:"email" gorm:"size:100;unique;not null" binding:"required,email,max=100"`
	Role     string `json:"role" gorm:"default:customer" binding:"oneof=customer admin"`
	// Orders   []Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	Cart Cart `json:"cart" gorm:"foreignKey:UserID"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
}
