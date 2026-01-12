package handlers

import "go.uber.org/zap"

type ServicesInterface interface {
	ProductService
	UserService
	CartService
}

type Handler struct {
	service ServicesInterface
	log     *zap.SugaredLogger
}

func NewHandler(service ServicesInterface, log *zap.SugaredLogger) *Handler {
	return &Handler{service: service, log: log}
}
