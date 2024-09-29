package main

import (
  "database/sql"
  "fmt"
  "log"
  "strconv"
  "github.com/gofiber/fiber/v2"
  _ "github.com/lib/pq"
)

const (
  host     = "localhost"  // or the Docker service name if running in another container
  port     = 5432         // default PostgreSQL port
  user     = "myuser"     // as defined in docker-compose.yml
  password = "mypassword" // as defined in docker-compose.yml
  dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *sql.DB

type Product struct {
	ID int `json:id`
	Name string `json:"name"`
	Price int `json:"price"`
}

func main() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  
  sdb, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    log.Fatal(err)
  }

  db = sdb
  defer db.Close()

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  app := fiber.New()

  app.Get("/product/:id", getProductHandler)
  app.Get("/products", getProductsHandler)
  app.Post("/product", createProductHandler)
  app.Put("/product/:id", updateProductHandler)
  app.Delete("product/:id", deleteProductHandler)

  app.Listen(":8080")
}

func getProduct(id int) (Product, error) {
  var p Product
  row := db.QueryRow(`SELECT id, name, price FROM products WHERE id = $1;`, id)
  err := row.Scan(&p.ID, &p.Name, &p.Price)
  if err != nil {
    return Product{}, err
  }
  return p, nil
}

func getProductHandler(c *fiber.Ctx) error {
  productId, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }
  product, err := getProduct(productId)
  if err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }
  return c.JSON(product)
}

func getProducts() ([]Product, error) {
  rows, err := db.Query("SELECT id, name, price FROM products")
  if err != nil {
    return nil, err
  }
  defer rows.Close()
  var products []Product
  for rows.Next() {
    var p Product
    err := rows.Scan(&p.ID, &p.Name, &p.Price)
    if err != nil {
      return nil, err
    }
    products = append(products, p)
  }
  if err = rows.Err(); err != nil {
    return nil, err
  }
  return products, err
}

func getProductsHandler(c *fiber.Ctx) error {
  products, err := getProducts()
  if err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }
  return c.JSON(products)
}

func createProduct(product *Product) error {
  _, err := db.Exec(
    "INSERT INTO products(name, price) VALUES ($1, $2)",
    product.Name,
    product.Price,
  )
  return err
}

func createProductHandler(c *fiber.Ctx) error {
  p := new(Product)
  if err := c.BodyParser(p); err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }

  err := createProduct(p)
  if err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }
  return c.JSON(p)
}

func updateProduct(id int, product *Product) (Product, error) {
  var p Product
  query := `
    UPDATE products
    SET name = $1, price = $2
    WHERE id = $3
    RETURNING id, name, price;
  `
  err := db.QueryRow(query, product.Name, product.Price, id).Scan(&p.ID, &p.Name, &p.Price)
  if err != nil {
    return Product{}, err
  }

  return p, nil
}

func updateProductHandler(c *fiber.Ctx) error {
  productId, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }

  p := new(Product)
  if err := c.BodyParser(p); err != nil {
    return c.SendStatus(fiber.StatusBadRequest)
  }

  fmt.Printf("Parsed Product: %+v\n", p)

  product, err := updateProduct(productId, p)
  if err != nil {
    return c.Status(fiber.StatusBadRequest).SendString("Product update failed.")
  }

  return c.JSON(product)
}

func deleteProduct(id int) error {
  _, err := db.Exec(
    "DELETE FROM products WHERE id=$1;",
    id,
  )
  return err
}

func deleteProductHandler(c *fiber.Ctx) error {
  productId, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(fiber.StatusBadRequest).SendString("Convert failed")
  }
  err = deleteProduct(productId)
  if err != nil {
    return c.Status(fiber.StatusBadRequest).SendString("Delete failed")
  }
  return c.SendStatus(fiber.StatusNoContent)
}