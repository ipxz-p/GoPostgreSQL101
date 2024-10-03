package main

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
  Data string `json:"data"`
}

type Pubsub struct {
  subs []chan Message
  mu sync.Mutex
}

func (ps *Pubsub) Subscribe() chan Message {
  ps.mu.Lock()
  defer ps.mu.Unlock()
  ch := make(chan Message, 1)
  ps.subs = append(ps.subs, ch)
  return ch
}

func (ps *Pubsub) Punlish(msg *Message) {
  ps.mu.Lock()
  defer ps.mu.Unlock()
  for _, sub := range ps.subs {
    sub <- *msg
  }
}

func main() {
  app := fiber.New()

  pubsub := &Pubsub{}

  app.Post("/publisher", func(c *fiber.Ctx) error {
    message := new(Message);
    if err := c.BodyParser(message); err != nil {
      return c.SendStatus(fiber.StatusBadRequest)
    }
    pubsub.Punlish(message)
    return c.JSON(fiber.Map{
      "message": "add to subscriber",
    }) 
  })
  
  sub := pubsub.Subscribe()
  go func ()  {
    for msg := range sub {
      fmt.Println("Receive message: ", msg)
    }
  }()

  app.Listen(":8080")
}