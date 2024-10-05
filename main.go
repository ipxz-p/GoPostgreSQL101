package main

import (
	"fmt"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/ipxz-p/GoPostgreSQL101/adapters"
	"github.com/ipxz-p/GoPostgreSQL101/entities"
	"github.com/ipxz-p/GoPostgreSQL101/usecases"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
  host     = "localhost"
  port     = 5432
  user     = "myuser"
  password = "mypassword"
  dbname   = "mydatabase"
)

func main() {
  app := fiber.New()

  dsn := fmt.Sprintf("host=%s port=%d user=%s "+
  "password=%s dbname=%s sslmode=disable",
  host, port, user, password, dbname)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
    panic("failed to connect to database")
  }

  if err := db.AutoMigrate(&entities.Order{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

  orderRepo := adapters.NewGormOrderRepository(db)
  orderService := usecases.NewOrderService(orderRepo)
  orderHandler := adapters.NewHttpOrderHandler(orderService)

  app.Post("/order", orderHandler.CreateOrder)
  app.Get("/order/:id", orderHandler.GetOrder)

  log.Fatal(app.Listen(":8000"))
}