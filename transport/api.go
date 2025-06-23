package transport

import (
	"context"
	"graphql-backend/app"
	"graphql-backend/graph/model"
	httptrans "graphql-backend/pkg/http-transport"
)

type API interface {
	CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error)
	UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error)
	Product(ctx context.Context, id string) (*model.Product, error)
	Products(ctx context.Context, limit *int32, offset *int32, category *string) ([]*model.Product, error)

	PlaceOrder(ctx context.Context, productIds []string) (*model.Order, error)
	Orders(ctx context.Context, limit *int32, offset *int32) ([]*model.Order, error)
	Order(ctx context.Context, id string) (*model.Order, error)

	Me(ctx context.Context) (*model.User, error)
	Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error)
}

type api struct {
	query   app.Query
	service app.Service
}

func (a api) Me(ctx context.Context) (*model.User, error) {
	userID := httptrans.GetUserFromContext(ctx).UserID
	user, err := a.query.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := UserRes{}
	res.Bind(user)

	return res.Res, nil
}

func (a api) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	result, err := a.service.Login(ctx, app.LoginParams{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	res := AuthPayloadRes{}
	res.Bind(result)

	return res.Res, nil
}

func (a api) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	product, err := a.service.CreateProduct(ctx, app.CreateProductParams{
		Name:        input.Name,
		Description: *input.Description,
		Price:       input.Price,
		InStock:     input.InStock,
		Category:    input.Category,
	})
	if err != nil {
		return nil, err
	}

	res := ProductRes{}
	res.Bind(product)

	return res.Res, nil
}

func (a api) UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error) {
	product, err := a.service.UpdateProduct(ctx, app.UpdateProductParams{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		InStock:     input.InStock,
		Category:    input.Category,
	})
	if err != nil {
		return nil, err
	}

	res := ProductRes{}
	res.Bind(product)

	return res.Res, nil
}

func (a api) Product(ctx context.Context, id string) (*model.Product, error) {
	product, err := a.query.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	res := ProductRes{}
	res.Bind(product)

	return res.Res, nil
}

func (a api) Products(ctx context.Context, limit *int32, offset *int32, category *string) ([]*model.Product, error) {
	es, err := a.query.GetProducts(ctx, app.ProductsParams{
		Limit:    limit,
		Offset:   offset,
		Category: category,
	})
	if err != nil {
		return nil, err
	}

	res := ProductsRes{}
	res.Bind(es)

	return res.Res, nil
}

func (a api) PlaceOrder(ctx context.Context, productIds []string) (*model.Order, error) {
	userID := httptrans.GetUserFromContext(ctx).UserID
	order, err := a.service.PlaceOrder(ctx, app.PlaceOrderParams{
		ProductIDs: productIds,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	res := OrderRes{}
	res.Bind(order)

	return res.Res, nil
}

func (a api) Orders(ctx context.Context, limit *int32, offset *int32) ([]*model.Order, error) {
	userID := httptrans.GetUserFromContext(ctx).UserID
	es, err := a.query.GetOrders(ctx, app.OrdersParams{
		Limit:  limit,
		Offset: offset,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	res := OrdersRes{}
	res.Bind(es)

	return res.Res, nil
}

func (a api) Order(ctx context.Context, id string) (*model.Order, error) {
	userID := httptrans.GetUserFromContext(ctx).UserID
	o, err := a.query.GetOrder(ctx, app.OrderParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	res := OrderRes{}
	res.Bind(o)

	return res.Res, nil
}

func NewAPI(query app.Query, service app.Service) API {
	return &api{
		query:   query,
		service: service,
	}
}
