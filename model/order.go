package model

// Order represents the structure of an order in the system
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
