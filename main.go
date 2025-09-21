package main

import (
	"context"
	"log"
	"net/http"
	"postgre_advanced/database"
	"postgre_advanced/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db := database.Connect()
	defer db.Close(context.Background())
	healthHandler := handlers.NewHealthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	r := mux.NewRouter()
	r.HandleFunc("/check", healthHandler.CheckServer).Methods("GET")
	r.HandleFunc("/db", healthHandler.CheckDatabase).Methods("GET")
	r.HandleFunc("/create-table", healthHandler.CreateUsersTable).Methods("POST")
	r.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
