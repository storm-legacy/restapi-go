package model

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// * DATABASE STRUCTURE
type User struct {
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// * CONFIGURATION
var secretKey = "somesecretjwt"

// Connect to database
func GetDbConnection() (db *gorm.DB, err error) {
	// Connect to database via dsn string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to database")

	fmt.Println("Migrating tables")
	err = db.AutoMigrate(User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
