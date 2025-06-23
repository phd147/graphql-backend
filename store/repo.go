package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"graphql-backend/app"
	"graphql-backend/entity"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type UserMap map[string]entity.User

type ProductMap map[string]entity.Product

type OrderMap map[string]entity.Order

// this repo implements the app.Repo interface
// we will use in-memory data for simplicity, and interval update it to json file
type repo struct {
	mu sync.RWMutex

	userMap    UserMap
	productMap ProductMap
	orderMap   OrderMap
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.userMap {
		if user.Email == email {
			return user, nil
		}
	}

	return entity.User{}, errors.New("user not found")
}

func (r *repo) CreateProduct(ctx context.Context, e entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.productMap[e.ID]; exists {
		return errors.New("product with the given ID already exists")
	}

	r.productMap[e.ID] = e
	return nil
}

func (r *repo) UpdateProduct(ctx context.Context, e entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.productMap[e.ID]; !exists {
		return errors.New("product not found")
	}

	r.productMap[e.ID] = e
	return nil
}

func (r *repo) CreateOrder(ctx context.Context, e entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orderMap[e.ID]; exists {
		return errors.New("order with the given ID already exists")
	}

	r.orderMap[e.ID] = e
	return nil
}

func (r *repo) GetProductsByIDs(ctx context.Context, ids []string) ([]entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]entity.Product, 0, len(ids))
	for _, id := range ids {
		if product, ok := r.productMap[id]; ok {
			products = append(products, product)
		}
	}
	return products, nil
}

func (r *repo) GetOrders(ctx context.Context, prs app.OrdersParams) ([]entity.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	limit := *prs.Limit
	offset := *prs.Offset

	var orders []entity.Order
	for _, order := range r.orderMap {
		if order.UserID == prs.UserID {
			orders = append(orders, order)
		}
	}

	start := offset
	end := offset + limit
	if int(start) > len(orders) {
		return []entity.Order{}, nil
	}
	if int(end) > len(orders) {
		end = int32(len(orders))
	}

	return orders[start:end], nil
}

func (r *repo) GetOrder(ctx context.Context, prs app.OrderParams) (entity.Order, error) {
	order, ok := r.orderMap[prs.ID]
	if !ok {
		return entity.Order{}, errors.New("order not found")
	}

	if prs.UserID != order.UserID {
		return entity.Order{}, errors.New("order not found")
	}

	return order, nil
}

func (r *repo) GetProductByID(ctx context.Context, id string) (entity.Product, error) {
	product, ok := r.productMap[id]
	if !ok {
		return entity.Product{}, errors.New("product not found")
	}

	return product, nil
}

func (r *repo) GetProducts(ctx context.Context, prs app.ProductsParams) ([]entity.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	limit := *prs.Limit
	offset := *prs.Offset

	var products []entity.Product
	for _, product := range r.productMap {
		if prs.Category == nil || *prs.Category == "" {
			products = append(products, product)
			continue
		}

		if strings.Contains(product.Category, *prs.Category) {
			products = append(products, product)
			continue
		}
	}

	start := offset
	end := offset + limit
	if int(start) > len(products) {
		return []entity.Product{}, nil
	}
	if int(end) > len(products) {
		end = int32(len(products))
	}

	return products[start:end], nil
}

func (r *repo) GetUsersByIDs(ctx context.Context, userIDs []string) ([]entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]entity.User, 0, len(userIDs))
	for _, id := range userIDs {
		if user, ok := r.userMap[id]; ok {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID string) (entity.User, error) {
	user, ok := r.userMap[userID]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}

	return user, nil
}

func loadMapFromFile[T any](filename string, out *map[string]T) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println("Failed to close file", filename, ":", err)
		}
	}()
	dec := json.NewDecoder(f)
	return dec.Decode(out)
}

func NewRepo(ctx context.Context) app.Repo {
	dir := filepath.Join("store", "data")
	usersPath := filepath.Join(dir, "users.json")
	productsPath := filepath.Join(dir, "products.json")
	ordersPath := filepath.Join(dir, "orders.json")

	userMap := UserMap{}
	productMap := ProductMap{}
	orderMap := OrderMap{}

	// Try to load from files, fallback to seed if not found
	_ = loadMapFromFile(usersPath, (*map[string]entity.User)(&userMap))
	_ = loadMapFromFile(productsPath, (*map[string]entity.Product)(&productMap))
	_ = loadMapFromFile(ordersPath, (*map[string]entity.Order)(&orderMap))

	// If userMap is empty, seed data for testing purposes
	if len(userMap) == 0 {
		adminID := uuid.NewString()
		customerID := uuid.NewString()
		userMap[adminID] = entity.User{
			ID:       adminID,
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "secret",
			Role:     "Admin",
		}
		userMap[customerID] = entity.User{
			ID:       customerID,
			Name:     "Customer User",
			Email:    "customer@example.com",
			Password: "secret",
			Role:     "Customer",
		}
	}

	r := &repo{
		mu:         sync.RWMutex{},
		userMap:    userMap,
		productMap: productMap,
		orderMap:   orderMap,
	}

	// write data to file in a separate goroutine and periodically update it
	go r.WriteDataToFile(ctx)

	return r
}

func (r *repo) WriteDataToFile(ctx context.Context) {
	dir := filepath.Join("store", "data")
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Failed to create data directory:", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	write := func(filename string, data any) {
		fpath := filepath.Join(dir, filename)
		f, err := os.Create(fpath)
		if err != nil {
			fmt.Println("Failed to create file", fpath, ":", err)
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				fmt.Println("Failed to close file", fpath, ":", err)
			}
		}(f)
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			fmt.Println("Failed to write data to", fpath, ":", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stop writing data to file due to context cancellation")
			return
		case <-ticker.C:
			r.mu.RLock()

			write("users.json", r.userMap)
			write("products.json", r.productMap)
			write("orders.json", r.orderMap)
			r.mu.RUnlock()
		}
	}
}
