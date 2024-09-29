package main

import (
	"fmt"
	"log"
  "os"
  "time"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "database/sql"
	// "fmt"
	// "log"
	// "strconv"
	"github.com/gofiber/fiber/v2"
  "github.com/joho/godotenv"
  "github.com/golang-jwt/jwt/v4"
	// _ "github.com/lib/pq"
)

const (
  host     = "localhost"  // or the Docker service name if running in another container
  port     = 5432         // default PostgreSQL port
  user     = "myuser"     // as defined in docker-compose.yml
  password = "mypassword" // as defined in docker-compose.yml
  dbname   = "mydatabase" // as defined in docker-compose.yml
)

func requireAuth(c *fiber.Ctx) error {
  cookie := c.Cookies("jwt")

  token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(os.Getenv("JWT_KEY")), nil
  })
  if err != nil || !token.Valid {
    return c.SendStatus(fiber.StatusUnauthorized)
  }

  return c.Next()
}

func main() {
  err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
  dsn := fmt.Sprintf("host=%s port=%d user=%s "+
  "password=%s dbname=%s sslmode=disable",
  host, port, user, password, dbname)

  newLogger := logger.New(
    log.New(os.Stdout, "\r\n", log.LstdFlags),
    logger.Config{
      SlowThreshold: time.Second,
      LogLevel:      logger.Info,
      Colorful:      true,
    },
  )

  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: newLogger,
  })

  if err != nil {
    panic("failed to connect to database")
  }

  db.AutoMigrate(&Book{}, &Publisher{}, &Author{}, &AuthorBook{})

  publisher := Publisher{
    Details: "Publisher Details",
    Name:    "Publisher Name",
  }
  _ = createPublisher(db, &publisher)

  author := Author{
    Name: "Author Name",
  }
  _ = createAuthor(db, &author)

  book := Book{
    Name:        "Book Title",
    Author:      "Book Author",
    Description: "Book Description",
    PublisherID: publisher.ID,
    Authors:     []Author{author},
  }
  _ = createBookWithAuthor(db, &book, []uint{author.ID})
  app := fiber.New()

  // app.Use("/books", requireAuth)
  // app.Get("/books", func(c *fiber.Ctx) error {
  //   return c.JSON(getBooks(db))
  // })
  // app.Post("/reg", func(c *fiber.Ctx) error {
  //   user := new(User)

  //   if err := c.BodyParser(user); err != nil {
  //     return c.SendStatus(fiber.StatusBadRequest)
  //   }

  //   err = createUser(db, user)

  //   if err != nil {
  //     return c.SendStatus(fiber.StatusBadRequest)
  //   }

  //   return c.JSON(fiber.Map{
  //     "message": "Reg success",
  //   })
  // })

  // app.Post("/login", func(c *fiber.Ctx) error {
  //   user := new(User)
  //   if err := c.BodyParser(user); err != nil {
  //     return c.SendStatus(fiber.StatusBadRequest)
  //   }
  //   token, err := loginUser(db, user)
  //   if err != nil {
  //     return c.SendStatus(fiber.StatusBadRequest)
  //   }

  //   c.Cookie(&fiber.Cookie{
  //     Name: "jwt",
  //     Value: token,
  //     Expires: time.Now().Add(time.Hour * 72),
  //     HTTPOnly: true,
  //   })

  //   return c.JSON(fiber.Map{
  //     "token": "Login success",
  //   })
  // })

  app.Listen(":8080")
}
