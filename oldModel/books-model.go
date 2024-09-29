package main

import (
	"fmt"
	"log"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name string `json:"name"`
	Author string `json:"author"`
	Description string `json:"description"`
	Price uint `json:"price"`
}

func createBook(db *gorm.DB, book *Book){
	result := db.Create(book)

	if result.Error != nil {
		log.Fatalf("Error creating book: %s", result.Error)
	}

	fmt.Println("Successful")
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return &book
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return books
}

func updateBook(db *gorm.DB, book *Book) {
	result := db.Save(&book)
	if result.Error != nil {
		log.Fatalf("Error updating book: %v", result.Error)
	  }
	  fmt.Println("Book updated successfully")
}

func DeleteBook(db *gorm.DB, id uint) {
	var book Book
	result := db.Delete(&book, id)
	if result.Error != nil {
	  log.Fatalf("Error deleting book: %v", result.Error)
	}
	fmt.Println("Book deleted successfully")
}

func SearchBook(db *gorm.DB, bookName string) []Book {
	var books []Book
	result := db.Where("name = ?", bookName).Order("price desc").Find(&books)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
	}
	return books
}