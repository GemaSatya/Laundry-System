package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GemaSatya/LaundrySystem/auth"
	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/utils"
	"github.com/joho/godotenv"
)

func main() {
	
	mux := http.NewServeMux()

	err := godotenv.Load()
	if err != nil{
		fmt.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")

	database.ConnectDatabase()
	fmt.Println("Start")

	// User Handler
	mux.HandleFunc("POST /users", utils.CreateUserHandler)
	mux.HandleFunc("POST /customers", utils.CreateCustomerHandler)
	mux.HandleFunc("POST /orders", utils.CreateOrderHandler)
	
	mux.HandleFunc("GET /customers", utils.GetAllCustomersHandler)
	mux.HandleFunc("GET /orders/me", utils.GetMyOrdersHandler)
	mux.HandleFunc("GET /order-details", utils.GetOrderByIDHandler)
	mux.HandleFunc("GET /orders", utils.GetAllOrdersHandler)

	// Authentication
	// mux.HandleFunc("/register", auth.RegisterUser)
	mux.HandleFunc("/login", auth.LoginUser)
	mux.HandleFunc("/logout", auth.LogoutUser)

	if err := http.ListenAndServe(":" + port, mux); err != nil{
		fmt.Println("Server started at port " + port)
		panic(err)
	}

}