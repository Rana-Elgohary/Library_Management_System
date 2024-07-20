package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Connect() {
	dsn := "root:rana5578652*@tcp(127.0.0.1:3306)/Library_Management_System_PyramakerzTask?charset=utf8mb4&parseTime=True&loc=Local"
	// &gorm.Config{} creates a new instance of gorm.Config
	// The gorm.Config struct allows you to configure various aspects of how GORM interacts with the database. It can include settings like logging preferences, naming conventions for tables
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	db = d
	fmt.Println("Successfully connected to the database")
}

func GetDB() *gorm.DB {
	return db
}
