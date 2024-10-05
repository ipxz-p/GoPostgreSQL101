package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/ipxz-p/GoPostgreSQL101/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderUseCase struct {
	mock.Mock
}

func (m *MockOrderUseCase) CreateOrder(order entities.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderUseCase) GetOrder(id uint) (*entities.Order, error) {
	args := m.Called(id)
	if order, ok := args.Get(0).(*entities.Order); ok {
		return order, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestHttpOrderHandler_CreateOrder(t *testing.T) {
	app := fiber.New()
	mockUseCase := new(MockOrderUseCase)
	handler := NewHttpOrderHandler(mockUseCase)

	app.Post("/order", handler.CreateOrder)

	t.Run("CreateOrder_Success", func(t *testing.T) {
		order := entities.Order{ID: 1, Total: 10}
		orderJSON, _ := json.Marshal(order)

		mockUseCase.On("CreateOrder", order).Return(nil)

		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(orderJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		var createdOrder entities.Order
		json.NewDecoder(res.Body).Decode(&createdOrder)
		assert.Equal(t, order, createdOrder)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("CreateOrder_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CreateOrder_TotalMustBePositive", func(t *testing.T) {
		order := entities.Order{ID: 1, Total: -5}
		orderJSON, _ := json.Marshal(order)

		mockUseCase.On("CreateOrder", order).Return(errors.New("Total must be positive"))

		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(orderJSON))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req, -1)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

		var errorResponse fiber.Map
		json.NewDecoder(res.Body).Decode(&errorResponse)
		assert.Equal(t, "Total must be positive", errorResponse["error"])

		mockUseCase.AssertExpectations(t)
	})
}