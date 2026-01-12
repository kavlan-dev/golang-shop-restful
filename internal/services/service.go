package services

type StorageInterface interface {
	ProductStorage
	UserStorage
	CartStorage
}

type Services struct {
	storage StorageInterface
}

func NewServices(db StorageInterface) *Services {
	return &Services{storage: db}
}
