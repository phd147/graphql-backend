package app

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"graphql-backend/entity"
	"graphql-backend/pkg/http-transport"
	"time"
)

const AccessTokenExpiration = 2 * time.Hour
const RefreshTokenExpiration = 7 * 24 * time.Hour

type Service interface {
	CreateProduct(ctx context.Context, prs CreateProductParams) (entity.Product, error)
	UpdateProduct(ctx context.Context, prs UpdateProductParams) (entity.Product, error)

	PlaceOrder(ctx context.Context, prs PlaceOrderParams) (entity.Order, error)
	Login(ctx context.Context, prs LoginParams) (LoginResult, error)
}

type Repo interface {
	GetOrders(ctx context.Context, prs OrdersParams) ([]entity.Order, error)
	GetOrder(ctx context.Context, prs OrderParams) (entity.Order, error)

	GetProductByID(ctx context.Context, id string) (entity.Product, error)
	GetProducts(ctx context.Context, prs ProductsParams) ([]entity.Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error)

	GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error)
	GetUserByID(ctx context.Context, id string) (entity.User, error)

	GetUserByEmail(ctx context.Context, email string) (entity.User, error)

	CreateProduct(ctx context.Context, e entity.Product) error
	UpdateProduct(ctx context.Context, e entity.Product) error

	CreateOrder(ctx context.Context, e entity.Order) error
}

type service struct {
	repo       Repo
	jwtHandler http_transport.JwtHandler
}

func (s service) Login(ctx context.Context, prs LoginParams) (LoginResult, error) {
	if prs.Email == "" || prs.Password == "" {
		return LoginResult{}, errors.New("email and password cannot be empty")
	}

	user, err := s.repo.GetUserByEmail(ctx, prs.Email)
	if err != nil {
		return LoginResult{}, err
	}

	if user.Password != prs.Password {
		return LoginResult{}, errors.New("invalid credentials")
	}

	accessToken, err := s.jwtHandler.GenerateToken(ctx, http_transport.UserClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "graphql-backend",
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Audience:  jwt.ClaimStrings{"graphql-ecommerce-client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
		},
	})
	if err != nil {
		return LoginResult{}, err
	}
	refreshToken, err := s.jwtHandler.GenerateToken(ctx, http_transport.UserClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "graphql-backend",
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			NotBefore: jwt.NewNumericDate(time.Now().UTC()),
			Audience:  jwt.ClaimStrings{"graphql-ecommerce-client"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
		},
	})
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,

		User: user,
	}, nil
}

func (s service) PlaceOrder(ctx context.Context, prs PlaceOrderParams) (entity.Order, error) {
	if len(prs.ProductIDs) == 0 {
		return entity.Order{}, errors.New("product IDs cannot be empty")
	}

	// Fetch products
	products, err := s.repo.GetProductsByIDs(ctx, prs.ProductIDs)
	if err != nil {
		return entity.Order{}, err
	}

	if len(products) == 0 {
		return entity.Order{}, errors.New("no products found for the given IDs")
	}

	// Create order
	var totalPrice float64
	for _, product := range products {
		totalPrice += product.Price
	}
	order := entity.Order{
		ID:         uuid.NewString(),
		UserID:     prs.UserID,
		ProductIDs: prs.ProductIDs,
		Total:      totalPrice,
		Status:     entity.OrderStatusPending,
		CreatedAt:  time.Now(),
	}

	err = s.repo.CreateOrder(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

func (s service) CreateProduct(ctx context.Context, prs CreateProductParams) (entity.Product, error) {
	product := entity.Product{
		ID:          uuid.NewString(),
		Name:        prs.Name,
		Description: prs.Description,
		Price:       prs.Price,
		InStock:     prs.InStock,
		Category:    prs.Category,
	}

	err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func (s service) UpdateProduct(ctx context.Context, prs UpdateProductParams) (entity.Product, error) {
	product, err := s.repo.GetProductByID(ctx, prs.ID)
	if err != nil {
		return entity.Product{}, err
	}

	prs.BindToProduct(&product)
	err = s.repo.UpdateProduct(ctx, product)
	if err != nil {
		return entity.Product{}, err
	}

	return product, nil
}

func NewService(repo Repo, jwtHandler http_transport.JwtHandler) Service {
	return &service{repo: repo, jwtHandler: jwtHandler}
}

type CreateProductParams struct {
	Name        string
	Description string
	Price       float64
	InStock     int32
	Category    string
}

type UpdateProductParams struct {
	ID string

	Name        *string
	Description *string
	Price       *float64
	InStock     *int32
	Category    *string
}

func (p *UpdateProductParams) BindToProduct(e *entity.Product) {
	if e == nil {
		return
	}
	if p.Name != nil {
		e.Name = *p.Name
	}
	if p.Description != nil {
		e.Description = *p.Description
	}
	if p.Price != nil {
		e.Price = *p.Price
	}
	if p.InStock != nil {
		e.InStock = *p.InStock
	}
	if p.Category != nil {
		e.Category = *p.Category
	}
}

type PlaceOrderParams struct {
	UserID     string
	ProductIDs []string
}

type LoginParams struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string

	User entity.User
}
