package utils

import (
	"encoding/json"
	"net/http"

	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/model"
)

type createUserRequest struct {
	UserId   int    `json:"user_id"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type createCustomerRequest struct {
	CustomerId int    `json:"customer_id"`
	Nama       string `json:"nama"`
	NoHp       string `json:"no_hp"`
	Alamat     string `json:"alamat"`
}

// Buat handler untuk membuat user baru (Admin / Owner)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Nama == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "nama, email, password, and role are required", http.StatusBadRequest)
		return
	}

	newUser := model.User{
		UserId:   req.UserId,
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// Buat handler untuk membuat customer baru (Customer)
func CreateCustomerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createCustomerRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Nama == "" || req.NoHp == "" || req.Alamat == "" {
		http.Error(w, "nama, no_hp, and alamat are required", http.StatusBadRequest)
		return
	}

	newCustomer := model.Customer{
		CustomerId: req.CustomerId,
		Nama:       req.Nama,
		No_hp:      req.NoHp,
		Alamat:     req.Alamat,
	}

	if err := database.DB.Create(&newCustomer).Error; err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCustomer)
}