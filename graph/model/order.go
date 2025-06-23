package model

type Order struct {
	ID         string     `json:"id"`
	Products   []*Product `json:"products"`
	ProductIDs []string   `json:"productIds"`
	Total      float64    `json:"total"`
	CreatedAt  string     `json:"createdAt"`
	Status     string     `json:"status"`
	User       *User      `json:"user"`
	UserID     string     `json:"userId"`
}
