package usecases

import (
	"testing"

	"github.com/ipxz-p/GoPostgreSQL101/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order entities.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) FindByID(id uint) (*entities.Order, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestCreateOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)
	mockRepo.On("Save", mock.Anything).Return(nil)
	tests := []struct {
		name      string
		order     entities.Order
		wantError bool
		wantMsg   string
	}{
		{
			name: "Success",
			order: entities.Order{
				ID:    1,
				Total: 10,
			},
			wantError: false,
		},
		{
			name: "Error Total Must Be Positive",
			order: entities.Order{
				ID:    3,
				Total: -5,
			},
			wantError: true,
			wantMsg:   "Total must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.CreateOrder(tt.order)

			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, tt.wantMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}