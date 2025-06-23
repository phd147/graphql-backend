package transport

import (
	"graphql-backend/app"
	"graphql-backend/entity"
	"graphql-backend/graph/model"
)

type ProductsRes struct {
	Res []*model.Product `json:"products"`
}

func (r *ProductsRes) Bind(es []entity.Product) {
	r.Res = make([]*model.Product, len(es))
	for i, e := range es {
		r.Res[i] = &model.Product{
			ID:          e.ID,
			Name:        e.Name,
			Description: StringP(e.Description),
			Price:       e.Price,
			Category:    e.Category,
			InStock:     e.InStock,
		}
	}
}

type UsersRes struct {
	Res []*model.User `json:"users"`
}

func (r *UsersRes) Bind(es []entity.User) {
	r.Res = make([]*model.User, len(es))
	for i, e := range es {
		r.Res[i] = &model.User{
			ID:    e.ID,
			Name:  e.Name,
			Email: e.Email,
		}
	}
}

type OrdersRes struct {
	Res []*model.Order `json:"orders"`
}

func (r *OrdersRes) Bind(es []entity.Order) {
	r.Res = make([]*model.Order, len(es))
	for i, e := range es {
		r.Res[i] = &model.Order{
			ID:         e.ID,
			ProductIDs: e.ProductIDs,
			Total:      e.Total,
			CreatedAt:  e.CreatedAt.String(),
			Status:     string(e.Status),
			UserID:     e.UserID,
		}
	}
}

type OrderRes struct {
	Res *model.Order `json:"order"`
}

func (r *OrderRes) Bind(e entity.Order) {
	r.Res = &model.Order{
		ID:         e.ID,
		ProductIDs: e.ProductIDs,
		Total:      e.Total,
		CreatedAt:  e.CreatedAt.String(),
		Status:     string(e.Status),
		UserID:     e.UserID,
	}
}

type ProductRes struct {
	Res *model.Product `json:"product"`
}

func (r *ProductRes) Bind(e entity.Product) {
	r.Res = &model.Product{
		ID:          e.ID,
		Name:        e.Name,
		Description: StringP(e.Description),
		Price:       e.Price,
		InStock:     e.InStock,
		Category:    e.Category,
	}
}

type UserRes struct {
	Res *model.User `json:"user"`
}

func (r *UserRes) Bind(e entity.User) {
	r.Res = &model.User{
		ID:    e.ID,
		Name:  e.Name,
		Email: e.Email,
	}
}

type AuthPayloadRes struct {
	Res *model.AuthPayload `json:"authPayload"`
}

func (r *AuthPayloadRes) Bind(e app.LoginResult) {
	r.Res = &model.AuthPayload{
		AccessToken:  e.AccessToken,
		RefreshToken: e.RefreshToken,
		User: &model.User{
			ID:    e.User.ID,
			Name:  e.User.Name,
			Email: e.User.Email,
		},
	}
}

func StringP(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
