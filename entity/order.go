package entity

import "time"

type Order struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	ProductIDs []string    `json:"product_ids"`
	Total      float64     `json:"total"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     OrderStatus `json:"status"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "Pending"
	OrderStatusCompleted OrderStatus = "Completed"
	OrderStatusCancelled OrderStatus = "Cancelled"
)
