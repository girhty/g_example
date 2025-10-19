package handlers

import (
	"g_example/db"
	"log"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect() (*gorm.DB, error) {
	sql_db, db_conn__err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if db_conn__err != nil {
		return nil, db_conn__err
	}
	return sql_db, nil
}

func GetAll(c fiber.Ctx) error {
	database, db_err := connect()
	if db_err != nil {
		log.Fatalf("failed to connect database : %s\n", db_err.Error())
		return db_err
	}
	var Products []db.Product
	result := database.Find(&Products)
	if result.Error != nil {
		log.Fatal(result.Error.Error())
		return result.Error
	}
	return c.JSON(Products)
}

func GetProduct_By_id(c fiber.Ctx) error {
	database, db_err := connect()
	if db_err != nil {
		log.Fatalf("failed to connect database : %s\n", db_err.Error())
		return db_err
	}
	id := c.Params("id")
	var Product db.Product
	res := database.First(&Product, id)
	if res.Error != nil {
		return res.Error
	}
	return c.JSON(Product)
}
