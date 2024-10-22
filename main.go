package main

import (
	
	"fmt"
	"log"
	"net/http"

	"github.com/YosaZiege/golang-react-todo/middleware"
	"github.com/YosaZiege/golang-react-todo/router" // Corrected import path
	"github.com/gorilla/handlers"
	// "github.com/YosaZiege/golang-react-todo/server/middleware"                       // Import the handlers package
)



func main() {
    middleware.InitDB()
    r := router.Router() // Call the Router function from the router package

    // Allow CORS for all origins
    cors := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}), // Allow all origins
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Allow specific methods
        handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
    )

    fmt.Println("Starting the server on Port 9000")
    log.Fatal(http.ListenAndServe(":9000", cors(r))) // Start the server on port 9000 with CORS enabled
}
