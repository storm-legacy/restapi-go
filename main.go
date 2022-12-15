package main

import (
	"fmt"
	"jwt-api/model"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	err error
	db  *gorm.DB
)

func init() {
	// Collect variables from .env
	err = godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	db, err = model.GetDbConnection()
}

func main() {
	fmt.Println("Hello! OwO")
}
