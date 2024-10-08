package usecases

import (
	"errors"

	"github.com/ipxz-p/GoPostgreSQL101/entities"
)

type OrderUseCase interface {
	CreateOrder(order entities.Order) error
	GetOrder(id uint) (*entities.Order, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) OrderUseCase {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(order entities.Order) error {
	if order.Total <= 0 {
		return errors.New("Total must be positive")
	}
	return s.repo.Save(order)
}

func (s *OrderService) GetOrder(id uint) (*entities.Order, error) {
	return s.repo.FindByID(id)
}