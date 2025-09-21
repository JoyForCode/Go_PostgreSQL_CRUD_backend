package handlers

import (
	"encoding/json"
	"postgre_advanced/models"
	"postgre_advanced/database"
	"net/http"
	"strconv"
	"github.com/jackc/pgx/v5"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	db *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserHandler {
	return &UserHandler{db:db}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req models.CreateUserRequest
	err:=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}
	err=database.CreateUser(h.db,req)
	if err!=nil{
		http.Error(w, `{"error":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message":"User created successfully!"})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err:=database.GetAllUsers(h.db)
	if err!=nil{
		http.Error(w, `{"error":"Failed to query users"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars:=mux.Vars(r)
	id, err:=strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, err:=database.GetUserByID(h.db,id)
	if err!=nil{
		if err==pgx.ErrNoRows{
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}
	http.Error(w, `{"error":"Failed to get user"}`, http.StatusInternalServerError)
	return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars:=mux.Vars(r)
	id, err:=strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	var req models.UpdateUserRequest
	err=json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	err=database.UpdateUser(h.db,id,req)
	if err!=nil{
		if err==pgx.ErrNoRows{
			http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to update user"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message":"User updated successfully!"})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars:=mux.Vars(r)
	id, err:=strconv.Atoi(vars["id"])
	if err!=nil{
		http.Error(w, `{"error":"Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	err=database.DeleteUser(h.db,id)
	if err!=nil{
		if err==pgx.ErrNoRows{
			http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Failed to delete user"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message":"User deleted successfully!"})
}
