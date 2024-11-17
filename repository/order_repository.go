package repository

import (
	"golang-orders-app/model"
)

// OrderRepository defines methods for interacting with the orders data.
type OrderRepository interface {
	CreateOrder(order *Order) (int, error)  // Method to create a new order
	GetUser(username string) (*User, error) // Method to get an user by username
	ListOrders(transferStatus, archive string, limit, page int, userid int) ([]OrderAll, int, error)
	CancelOrder(consignmentID int) error
}

// Order represents an order in the repository layer.
type Order struct {
	ID                 int     `json:"id"`
	UserID             int     `json:"user_id"`
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
	OrderTypeID        int     `json:"order_type_id" validate:"required"`
	TotalFee           float64 `json:"total_fee"`
	CODFee             float64 `json:"cod_fee"`
	PromoDiscount      float64 `json:"promo_discount"`
	Discount           float64 `json:"discount"`
	DeliveryFee        float64 `json:"delivery_fee"`
	Archive            bool    `json:"archive"`
}

// OrderAll represents an order response in the repository layer.
type OrderAll struct {
	OrderConsignmentID string  `json:"order_consignment_id"`
	OrderCreatedAt     string  `json:"order_created_at"`
	OrderDescription   string  `json:"order_description"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name"`
	RecipientAddress   string  `json:"recipient_address"`
	RecipientPhone     string  `json:"recipient_phone"`
	OrderAmount        float64 `json:"order_amount"`
	DeliveryFee        float64 `json:"delivery_fee"`
	CODFee             float64 `json:"cod_fee"`
	PromoDiscount      float64 `json:"promo_discount"`
	Discount           float64 `json:"discount"`
	OrderStatus        string  `json:"order_status"`
	OrderType          string  `json:"order_type"`
	ItemType           string  `json:"item_type"`
	Instruction        string  `json:"instruction,omitempty"`
	TotalFee           float64 `json:"total_fee"`
}

// NewOrderFromModel converts a model.Order to repository.Order
func NewOrderFromModel(m *model.Order) *Order {
	return &Order{
		ID:                 m.ID,
		UserID:             m.UserID,
		StoreID:            m.StoreID,
		MerchantOrderID:    m.MerchantOrderID,
		RecipientName:      m.RecipientName,
		RecipientPhone:     m.RecipientPhone,
		RecipientAddress:   m.RecipientAddress,
		RecipientCity:      m.RecipientCity,
		RecipientZone:      m.RecipientZone,
		RecipientArea:      m.RecipientArea,
		DeliveryType:       m.DeliveryType,
		ItemType:           m.ItemType,
		SpecialInstruction: m.SpecialInstruction,
		ItemQuantity:       m.ItemQuantity,
		ItemWeight:         m.ItemWeight,
		AmountToCollect:    m.AmountToCollect,
		ItemDescription:    m.ItemDescription,
		OrderTypeID:        1,
		TotalFee:           m.TotalFee,
		CODFee:             m.CODFee,
		PromoDiscount:      0.00,
		Discount:           0.00,
		DeliveryFee:        m.DeliveryFee,
		Archive:            false,
	}
}
