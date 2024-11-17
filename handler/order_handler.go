package handler

import (
	"encoding/json"
	"fmt"
	"golang-orders-app/model"
	"golang-orders-app/repository"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
)

// OrderHandler struct holds the repository for the orders
type OrderHandler struct {
	orderRepo repository.OrderRepository
}

// NewOrderHandler initializes the OrderHandler
func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

// CreateOrder handles the POST request for creating an order
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate JWT Token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

	// Parse the JWT token to extract user ID
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the JWT secret key for verification
		return []byte("secret"), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from the token
	claims, _ := token.Claims.(jwt.MapClaims)
	userName := claims["username"].(string)
	user, err := h.orderRepo.GetUser(userName)

	userID := user.ID

	// Step 2: Parse and Validate Request Body
	var orderRequest struct {
		StoreID            int     `json:"store_id"`
		MerchantOrderID    string  `json:"merchant_order_id"`
		RecipientName      string  `json:"recipient_name"`
		RecipientPhone     string  `json:"recipient_phone"`
		RecipientAddress   string  `json:"recipient_address"`
		RecipientCity      int     `json:"recipient_city"`
		RecipientZone      int     `json:"recipient_zone"`
		RecipientArea      int     `json:"recipient_area"`
		DeliveryType       int     `json:"delivery_type"`
		ItemType           int     `json:"item_type"`
		SpecialInstruction string  `json:"special_instruction"`
		ItemQuantity       int     `json:"item_quantity"`
		ItemWeight         float64 `json:"item_weight"`
		AmountToCollect    float64 `json:"amount_to_collect"`
		ItemDescription    string  `json:"item_description"`
	}

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Step 3: Validate Required Fields
	errors := make(map[string][]string)

	if orderRequest.StoreID == 0 {
		errors["store_id"] = append(errors["store_id"], "The store field is required")
	}

	if orderRequest.RecipientName == "" {
		errors["recipient_name"] = append(errors["recipient_name"], "The recipient name field is required")
	}

	phoneRegex := regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`) // BD Number Validation
	if !phoneRegex.MatchString(orderRequest.RecipientPhone) {
		errors["recipient_phone"] = append(errors["recipient_phone"], "Invalid phone number")
	}

	if orderRequest.RecipientAddress == "" {
		errors["recipient_address"] = append(errors["recipient_address"], "The recipient address field is required")
	}

	if orderRequest.DeliveryType == 0 {
		errors["delivery_type"] = append(errors["delivery_type"], "The delivery type field is required")
	}

	if orderRequest.AmountToCollect == 0 {
		errors["amount_to_collect"] = append(errors["amount_to_collect"], "The amount to collect field is required")
	}

	if orderRequest.ItemQuantity == 0 {
		errors["item_quantity"] = append(errors["item_quantity"], "The item quantity field is required")
	}

	if orderRequest.ItemWeight == 0 {
		errors["item_weight"] = append(errors["item_weight"], "The item weight field is required")
	}

	if orderRequest.ItemType == 0 {
		errors["item_type"] = append(errors["item_type"], "The item type field is required")
	}

	if len(errors) > 0 {
		// Return validation errors
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Please fix the given errors",
			"type":    "error",
			"code":    422,
			"errors":  errors,
		})
		return
	}

	// Calculate Delivery Fee
	var deliveryFee float64
	switch {
	case orderRequest.RecipientCity == 1 && orderRequest.ItemWeight <= 0.5:
		deliveryFee = 60
	case orderRequest.RecipientCity == 1 && orderRequest.ItemWeight > 0.5 && orderRequest.ItemWeight <= 1:
		deliveryFee = 70
	case orderRequest.RecipientCity == 1 && orderRequest.ItemWeight > 1:
		deliveryFee = 70 + 15*math.Ceil(orderRequest.ItemWeight-1)
	default:
		// If RecipientCity != 1
		deliveryFee = 100 + 15*math.Ceil(orderRequest.ItemWeight-1)
	}

	// Calculate COD Fee (1% of AmountToCollect)
	codFee := orderRequest.AmountToCollect * 0.01

	totalFee := orderRequest.AmountToCollect + codFee + deliveryFee

	// Step 4: Create the Order
	order := model.Order{
		UserID:             int(userID),
		StoreID:            orderRequest.StoreID,
		MerchantOrderID:    orderRequest.MerchantOrderID,
		RecipientName:      orderRequest.RecipientName,
		RecipientPhone:     orderRequest.RecipientPhone,
		RecipientAddress:   orderRequest.RecipientAddress,
		RecipientCity:      orderRequest.RecipientCity,
		RecipientZone:      orderRequest.RecipientZone,
		RecipientArea:      orderRequest.RecipientArea,
		DeliveryType:       orderRequest.DeliveryType,
		ItemType:           orderRequest.ItemType,
		SpecialInstruction: orderRequest.SpecialInstruction,
		ItemQuantity:       orderRequest.ItemQuantity,
		ItemWeight:         orderRequest.ItemWeight,
		AmountToCollect:    orderRequest.AmountToCollect,
		ItemDescription:    orderRequest.ItemDescription,
		OrderTypeID:        1,
		TotalFee:           totalFee, // Optional field
		CODFee:             codFee,   // Optional field
		PromoDiscount:      0.00,     // Optional field
		Discount:           0.00,     // Optional field
		DeliveryFee:        deliveryFee,
		Archive:            false,
	}

	repoOrder := repository.NewOrderFromModel(&order) // Convert to repository order

	// Insert the order into the database
	consignmentID, err := h.orderRepo.CreateOrder(repoOrder) // Now passing the repository order
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 5: Respond with success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order Created Successfully",
		"type":    "success",
		"code":    200,
		"data": map[string]interface{}{
			"consignment_id":    consignmentID,
			"merchant_order_id": orderRequest.MerchantOrderID,
			"order_status":      "Pending",
			"delivery_fee":      deliveryFee,
		},
	})
}

// ListOrders handles the GET request for listing orders with pagination and filters.
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate JWT Token
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract user ID from the token
	claims, _ := token.Claims.(jwt.MapClaims)
	userName := claims["username"].(string)
	user, err := h.orderRepo.GetUser(userName)
	
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userid := user.ID

	// Step 2: Extract and validate query parameters
	query := r.URL.Query()

	transferStatus := query.Get("transfer_status")
	if transferStatus == "1" {
		transferStatus = "Pending"
	} else {
		transferStatus = "Cancel"
	}
	archive := query.Get("archive")
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default page
	}

	// Step 3: Call the repository to fetch orders
	orders, total, err := h.orderRepo.ListOrders(transferStatus, archive, limit, page, userid)
	if err != nil {
		log.Printf("Failed to fetch orders: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Prepare pagination details
	lastPage := (total + limit - 1) / limit
	totalInPage := len(orders)

	// Step 5: Respond with the order data
	response := map[string]interface{}{
		"message": "Orders successfully fetched.",
		"type":    "success",
		"code":    200,
		"data": map[string]interface{}{
			"data":          orders,
			"total":         total,
			"current_page":  page,
			"per_page":      limit,
			"total_in_page": totalInPage,
			"last_page":     lastPage,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CancelOrderHandler handles the cancellation of an order
func (h *OrderHandler) CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the consignment ID from the URL path
	consignmentIDStr := chi.URLParam(r, "consignmentID")
	consignmentID, err := strconv.Atoi(consignmentIDStr)
	if err != nil {
		http.Error(w, `{"message": "Invalid consignment ID", "type": "error", "code": 400}`, http.StatusBadRequest)
		return
	}

	// Validate the JWT
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenString = tokenString[len("Bearer "):] // Remove "Bearer " prefix

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Cancel the order using the repository
	err = h.orderRepo.CancelOrder(consignmentID) // Use your repository instance here
	if err != nil {
		if err.Error() == "order already cancelled or not found" {
			http.Error(w, `{"message": "Please contact cx to cancel order", "type": "error", "code": 400}`, http.StatusBadRequest)
			return
		}

		http.Error(w, `{"message": "Internal server error", "type": "error", "code": 500}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message": "Order Cancelled Successfully",
		"type":    "success",
		"code":    200,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
