package main

import (
	"fmt"
	"log"
	"os"
	"bioskop/db"
	"bioskop/handler"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.InitDB() // Inisialisasi koneksi database
	defer db.CloseDB()

	app := fiber.New()

	app.Post("/users", handler.CreateUser)       
	app.Get("/users", handler.GetAllUsers)       
	app.Get("/users/:id", handler.GetUserByID)   
	app.Put("/users/:id", handler.UpdateUser)    
	app.Delete("/users/:id", handler.DeleteUser) 

	app.Post("/cities", handler.CreateCity)       
	app.Get("/cities", handler.GetAllCities)       
	app.Get("/cities/:id", handler.GetCityByID)   
	app.Put("/cities/:id", handler.UpdateCity)    
	app.Delete("/cities/:id", handler.DeleteCity) 

	app.Post("/cinemas", handler.CreateCinema)       
	app.Get("/cinemas", handler.GetAllCinemas)       
	app.Get("/cinemas/:id", handler.GetCinemaByID)   
	app.Put("/cinemas/:id", handler.UpdateCinema)    
	app.Delete("/cinemas/:id", handler.DeleteCinema) 

	fmt.Println("Server berjalan di http://"+os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"))
	log.Fatal(app.Listen(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"))) // Gunakan PORT dari .env atau default 3000
}