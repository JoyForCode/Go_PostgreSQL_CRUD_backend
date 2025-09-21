package handlers

import (
	"fmt"
	"net/http"
	"postgre_advanced/database"

	"github.com/jackc/pgx/v5"
)

type HealthHandler struct {
	db *pgx.Conn
}

func NewHealthHandler(db *pgx.Conn) *HealthHandler {
	return &HealthHandler{db:db}
}

func (h *HealthHandler) CheckServer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "Server is up and running"}`)
}

func (h *HealthHandler) CheckDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err:=database.TestConnection(h.db)
	if err!=nil{
		http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `{"status":"Database is working!"}`)
}

func (h *HealthHandler) CreateUsersTable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err:=database.CreateUsersTable(h.db)
	if err!=nil{
		http.Error(w, `{"error":"Failed to create users table"}`, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, `{"status":"Users table created successfully"}`)
}
