package adapters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipxz-p/GoPostgreSQL101/usecases"
	"github.com/ipxz-p/GoPostgreSQL101/entities"
)

type HttpOrderHandler struct {
	orderUseCase usecases.OrderUseCase
}

func NewHttpOrderHandler(useCase usecases.OrderUseCase) *HttpOrderHandler {
	return &HttpOrderHandler{orderUseCase: useCase}
}

func (h *HttpOrderHandler) CreateOrder(c *fiber.Ctx) error {
	var order entities.Order
	if err := c.BodyParser(&order); err != nil {
	  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.orderUseCase.CreateOrder(order); err != nil {
	  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
  
	return c.Status(fiber.StatusCreated).JSON(order)
  }

func (h *HttpOrderHandler) GetOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err !=  nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	order, err := h.orderUseCase.GetOrder(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}
	return c.Status(fiber.StatusCreated).JSON(order)
}