package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

// OrderRepositoryImpl is the struct that implements the OrderRepository interface.
type OrderRepositoryImpl struct {
	DB *sql.DB
}

// NewOrderRepository creates a new instance of OrderRepository.
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{DB: db}
}

// CreateOrder creates a new order in the database.
func (r *OrderRepositoryImpl) CreateOrder(order *Order) (int, error) {
	query := `INSERT INTO orders (userid, store_id, merchant_order_id, recipient_name, recipient_phone, 
    recipient_address, recipient_city, recipient_zone, recipient_area, delivery_type, item_type, 
    special_instruction, item_quantity, item_weight, amount_to_collect, item_description, order_type_id, 
    total_fee, cod_fee, promo_discount, discount, delivery_fee, archive) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23) 
    RETURNING id`
	fmt.Println("order fee", order.TotalFee)

	var consignmentID int
	err := r.DB.QueryRow(query,
		order.UserID, order.StoreID, order.MerchantOrderID, order.RecipientName,
		order.RecipientPhone, order.RecipientAddress, order.RecipientCity, order.RecipientZone,
		order.RecipientArea, order.DeliveryType, order.ItemType, order.SpecialInstruction,
		order.ItemQuantity, order.ItemWeight, order.AmountToCollect, order.ItemDescription,
		order.OrderTypeID, order.TotalFee, order.CODFee, order.PromoDiscount, order.Discount,
		order.DeliveryFee, order.Archive,
	).Scan(&consignmentID)

	if err != nil {
		return 0, fmt.Errorf("error creating order: %v", err)
	}
	return consignmentID, nil
}

// GetUser fetches a user from the database by username.
func (r *OrderRepositoryImpl) GetUser(username string) (*User, error) {
	query := `SELECT id FROM users WHERE username = $1`
	row := r.DB.QueryRow(query, username)

	var user User
	if err := row.Scan(&user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

// ListOrders fetches a list of orders from the database based on the given parameters.
func (r *OrderRepositoryImpl) ListOrders(transferStatus, archive string, limit, page int, userid int) ([]OrderAll, int, error) {
	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Query to fetch orders
	query :=
		`SELECT                
    o.id AS order_consignment_id,
    o.created_at AS order_created_at,
    o.item_description AS order_description,
    o.merchant_order_id,
    o.recipient_name,
    o.recipient_address,
    o.recipient_phone,
    o.amount_to_collect AS order_amount,
    o.delivery_fee,
    o.cod_fee,
    o.promo_discount,
    o.discount,
    o.order_status,
    o.order_type_id,
    o.item_type,
    o.special_instruction AS instruction,
    o.total_fee
FROM orders o
JOIN users u ON o.userid = $5
WHERE o.order_status = $1 AND o.archive = $2 AND o.userId = $5
LIMIT $3 OFFSET $4;
`

	rows, err := r.DB.Query(query, transferStatus, archive, limit, offset, userid)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching orders: %v", err)
	}
	defer rows.Close()

	var orders []OrderAll
	for rows.Next() {
		var order OrderAll
		if err := rows.Scan(
			&order.OrderConsignmentID,
			&order.OrderCreatedAt,
			&order.OrderDescription,
			&order.MerchantOrderID,
			&order.RecipientName,
			&order.RecipientAddress,
			&order.RecipientPhone,
			&order.OrderAmount,
			&order.DeliveryFee,
			&order.CODFee,
			&order.PromoDiscount,
			&order.Discount,
			&order.OrderStatus,
			&order.OrderType,
			&order.ItemType,
			&order.Instruction,
			&order.TotalFee,
		); err != nil {
			return nil, 0, fmt.Errorf("error scanning order: %v", err)
		}
		orders = append(orders, order)
	}

	// Count total orders for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM orders WHERE  order_status= $1 AND archive = $2`
	if err := r.DB.QueryRow(countQuery, transferStatus, archive).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("error counting orders: %v", err)
	}

	return orders, total, nil
}

// CancelOrder sets the order status to "Cancelled" for the given consignment ID.
func (r *OrderRepositoryImpl) CancelOrder(consignmentID int) error {
	query := `UPDATE orders SET order_status = 'Cancelled' WHERE id = $1 AND order_status != 'Cancelled'`
	result, err := r.DB.Exec(query, consignmentID)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to fetch affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("order already cancelled or not found")
	}

	return nil
}
