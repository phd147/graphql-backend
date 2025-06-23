package order

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestPlaceOrder(t *testing.T) {
	adminToken := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// Create a product to order
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "OrderProduct",
		"price":       20.0,
		"inStock":     10,
		"description": "desc",
		"category":    "OrderCat",
	}
	createReq.Var("input", input)
	tests.AuthRequest(createReq, adminToken)
	var createResp struct {
		CreateProduct struct{ ID string }
	}
	err := client.Run(context.TODO(), createReq, &createResp)
	require.NoError(t, err)
	productID := createResp.CreateProduct.ID

	// Place order as customer
	customerToken := tests.Login(t, tests.CustomerEmail, tests.CustomerPassword)
	orderReq := graphql.NewRequest(`mutation($ids: [ID!]!) { placeOrder(productIds: $ids) { id total status products { id name } user { id name } } }`)
	orderReq.Var("ids", []string{productID})
	tests.AuthRequest(orderReq, customerToken)
	var orderResp struct {
		PlaceOrder struct {
			ID       string
			Total    float64
			Status   string
			Products []struct {
				ID   string
				Name string
			}
			User struct {
				ID   string
				Name string
			}
		}
	}
	err = client.Run(context.TODO(), orderReq, &orderResp)
	require.NoError(t, err)
	require.Equal(t, 1, len(orderResp.PlaceOrder.Products))
}
