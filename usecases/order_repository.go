package usecases

import (
	"github.com/ipxz-p/GoPostgreSQL101/entities"
)

type OrderRepository interface {
	Save(order entities.Order) error
	FindByID(id uint) (*entities.Order, error)
}

