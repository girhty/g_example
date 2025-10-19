package main

import (
	"errors"
	"fmt"
	"g_example/db"
	"g_example/handlers"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	if _, err := os.Stat("./test.db"); errors.Is(err, os.ErrNotExist) {
		fmt.Print("database does not exist , building new db")
		if err := db.Create_db(); err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Print("db exists")
	}
	app := fiber.New()
	app.Get("/", handlers.GetAll)
	app.Get("/p/:id", handlers.GetProduct_By_id)

	log.Fatal(app.Listen(":3000"))
}
