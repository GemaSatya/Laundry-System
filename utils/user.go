package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/model"
)

type createUserRequest struct {
	Nama     string `json:"nama"`
	Username string `json:"username"`
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

func decodeCreateUserRequest(r *http.Request) (createUserRequest, error) {
	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		return createUserRequest{}, err
	}

	if len(raw) == 0 {
		return createUserRequest{}, fmt.Errorf("empty request body")
	}

	if raw[0] == '"' {
		var encoded string
		if err := json.Unmarshal(raw, &encoded); err != nil {
			return createUserRequest{}, err
		}

		raw = []byte(encoded)
	}

	var req createUserRequest
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		return createUserRequest{}, err
	}

	return req, nil
}

// Buat handler untuk membuat user baru (Admin / Owner)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := decodeCreateUserRequest(r)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	if req.Nama == "" || req.Username == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		http.Error(w, "nama, username, email, password, and role are required", http.StatusBadRequest)
		return
	}

	newUser := model.User{
		Nama:     req.Nama,
		Username: req.Username,
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