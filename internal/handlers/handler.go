package handlers

import "golang-shop-restful/internal/services"

type HandlersInterface interface {
	ProductHandler
	UserHandler
	CartHandler
}

type Handler struct {
	service services.ServicesInterface
}

func NewHandler(service services.ServicesInterface) *Handler {
	return &Handler{service: service}
}
