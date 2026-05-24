package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/model"
)

func requestRole(r *http.Request) string {
	role := r.Header.Get("X-User-Role")
	if role == "" {
		role = r.Header.Get("Role")
	}

	return strings.ToLower(strings.TrimSpace(role))
}

func GetAllCustomersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if requestRole(r) != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var customers []model.Customer
	if err := database.DB.Find(&customers).Error; err != nil {
		http.Error(w, "Failed to fetch customers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}