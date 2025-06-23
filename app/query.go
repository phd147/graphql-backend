package app

import (
	"context"
	"graphql-backend/entity"
)

type Query interface {
	GetProducts(ctx context.Context, prs ProductsParams) ([]entity.Product, error)
	GetProduct(ctx context.Context, id string) (entity.Product, error)

	GetOrders(ctx context.Context, prs OrdersParams) ([]entity.Order, error)
	GetOrder(ctx context.Context, prs OrderParams) (entity.Order, error)

	GetUser(ctx context.Context, id string) (entity.User, error)
}

type query struct {
	repo Repo
}

func (q *query) GetUser(ctx context.Context, id string) (entity.User, error) {
	return q.repo.GetUserByID(ctx, id)
}

func (q *query) GetOrders(ctx context.Context, prs OrdersParams) ([]entity.Order, error) {
	prs.SetDefaults()
	return q.repo.GetOrders(ctx, prs)
}

func (q *query) GetOrder(ctx context.Context, prs OrderParams) (entity.Order, error) {
	return q.repo.GetOrder(ctx, prs)
}

func (q *query) GetProduct(ctx context.Context, id string) (entity.Product, error) {
	return q.repo.GetProductByID(ctx, id)
}

func (q *query) GetProducts(ctx context.Context, prs ProductsParams) ([]entity.Product, error) {
	prs.SetDefaults()
	return q.repo.GetProducts(ctx, prs)
}

func NewQuery(repo Repo) Query {
	return &query{repo: repo}
}

type ProductsParams struct {
	Limit    *int32
	Offset   *int32
	Category *string
}

func (p *ProductsParams) SetDefaults() {
	if p.Limit == nil || *p.Limit <= 0 {
		defaultLimit := int32(10)
		p.Limit = &defaultLimit
	}
	if p.Offset == nil || *p.Offset < 0 {
		defaultOffset := int32(0)
		p.Offset = &defaultOffset
	}
	if p.Category == nil {
		defaultCategory := ""
		p.Category = &defaultCategory
	}
}

type OrdersParams struct {
	Limit  *int32
	Offset *int32
	UserID string
}

func (o *OrdersParams) SetDefaults() {
	if o.Limit == nil || *o.Limit <= 0 {
		defaultLimit := int32(10)
		o.Limit = &defaultLimit
	}
	if o.Offset == nil || *o.Offset < 0 {
		defaultOffset := int32(0)
		o.Offset = &defaultOffset
	}
}

type OrderParams struct {
	ID     string
	UserID string
}
