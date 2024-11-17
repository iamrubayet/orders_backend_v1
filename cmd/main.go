package main

import (
	"log"
	"net/http"

	"golang-orders-app/config"
	"golang-orders-app/handler"

	"golang-orders-app/repository"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories and handlers
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db) 
	userHandler := handler.NewUserHandler(userRepo)
	orderHandler := handler.NewOrderHandler(orderRepo)

	// Initialize Chi router
	r := chi.NewRouter()
	// Register routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", userHandler.LoginHandler)
		r.Post("/logout", userHandler.LogoutHandler)
		r.Post("/orders", orderHandler.CreateOrder)
		r.Get("/orders/all", orderHandler.ListOrders)
		r.Put("/orders/{consignmentID}/cancel", orderHandler.CancelOrderHandler)

	})

	// Start the HTTP server
	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
