package adapters

import (
	"gorm.io/gorm"
	"github.com/ipxz-p/GoPostgreSQL101/usecases"
	"github.com/ipxz-p/GoPostgreSQL101/entities"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) usecases.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) Save(order entities.Order) error {
	return r.db.Create(&order).Error
}

func (r *GormOrderRepository) FindByID(id uint) (*entities.Order, error) {
	var order entities.Order
	if err := r.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}