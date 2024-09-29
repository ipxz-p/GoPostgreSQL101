package main

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Password string
}

func createUser(db *gorm.DB, user *User) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashPassword)

	result := db.Create(user)
	if result.Error != nil {
		return err
	}

	return nil
}

func loginUser(db *gorm.DB, user *User) (string, error) {
	findUser := new(User)

	result := db.Where("email = ?", user.Email).First(findUser)
	if result.Error != nil {
		return "", result.Error
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(findUser.Password), 
		[]byte(user.Password),
	)
	if err != nil {
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = findUser.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return t, nil
}