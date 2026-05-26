package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/GemaSatya/LaundrySystem/database"
	"github.com/GemaSatya/LaundrySystem/model"
	"gorm.io/gorm"
)

type createOrderRequest struct {
	OrderId       int    `json:"order_id"`
	TanggalMasuk  string `json:"tanggal_masuk"`
	TanggalKeluar string `json:"tanggal_keluar"`
	TotalHarga    int    `json:"total_harga"`
	OrderDetail   []createOrderDetailRequest `json:"order_detail"`
}

type createOrderDetailRequest struct {
	JenisLayanan string `json:"jenis_layanan"`
	Berat        int    `json:"berat"`
	HargaPerKg   int    `json:"harga_per_kg"`
	Subtotal     int    `json:"subtotal"`
}

// GetOrderDetails fetches a single order with its details from the database
func GetOrderDetails(orderID int) (*model.Order, error) {
	var order model.Order
	
	result := database.DB.Preload("OrderDetail").First(&order, orderID)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &order, nil
}

// GetAllOrders fetches all orders from the database
func GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	
	result := database.DB.Preload("OrderDetail").Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	
	return orders, nil
}

// GetOrdersByCustomerID fetches all orders that belong to a specific customer.
func GetOrdersByCustomerID(customerID int) ([]model.Order, error) {
	var orders []model.Order

	result := database.DB.Preload("OrderDetail").Where("order_id = ?", customerID).Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

// GetOrderByIDHandler is an HTTP handler to fetch order details by ID
func GetOrderByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get order ID from query parameter
	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(orderID)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := GetOrderDetails(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// GetAllOrdersHandler is an HTTP handler to fetch all orders
func GetAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orders, err := GetAllOrders()
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetMyOrdersHandler returns orders for a single customer.
func GetMyOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customerIDValue := r.URL.Query().Get("customer_id")
	if customerIDValue == "" {
		customerIDValue = r.Header.Get("X-Customer-Id")
	}
	if customerIDValue == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	customerID, err := strconv.Atoi(customerIDValue)
	if err != nil || customerID <= 0 {
		http.Error(w, "Invalid customer_id", http.StatusBadRequest)
		return
	}

	var customer model.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	orders, err := GetOrdersByCustomerID(customerID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// CreateOrderHandler creates a new order from a JSON request body.
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.TanggalMasuk == "" || req.TanggalKeluar == "" {
		http.Error(w, "tanggal_masuk and tanggal_keluar are required", http.StatusBadRequest)
		return
	}

	if req.OrderId <= 0 {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	var customer model.Customer
	if err := database.DB.First(&customer, req.OrderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Customer not found for order_id", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to validate customer", http.StatusInternalServerError)
		return
	}

	orderDetails := make([]model.OrderDetail, 0, len(req.OrderDetail))
	for _, item := range req.OrderDetail {
		if item.JenisLayanan == "" {
			http.Error(w, "jenis_layanan is required for each order_detail item", http.StatusBadRequest)
			return
		}

		orderDetails = append(orderDetails, model.OrderDetail{
			Jenis_layanan: item.JenisLayanan,
			Berat:         item.Berat,
			Harga_per_kg:  item.HargaPerKg,
			Subtotal:      item.Subtotal,
		})
	}

	order := model.Order{
		OrderId:       req.OrderId,
		Tanggal_masuk:  req.TanggalMasuk,
		Tanggal_keluar: req.TanggalKeluar,
		Total_harga:    req.TotalHarga,
		OrderDetail:    orderDetails,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	/* var orderRequest struct{
		OrderId int `json:"order_id"`
		TanggalMasuk  string `json:"tanggal_masuk"`
		TanggalKeluar string `json:"tanggal_keluar"`
		TotalHarga    int    `json:"total_harga"`
	}

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	newOrder := model.Order{
		Tanggal_masuk: orderRequest.TanggalMasuk,
		Tanggal_keluar: orderRequest.TanggalKeluar,
		Total_harga: orderRequest.TotalHarga,
	}

	if err := database.DB.Create(&newOrder).Error; err != nil{
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	} */

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
