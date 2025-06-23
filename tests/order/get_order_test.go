package order

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestGetOrder(t *testing.T) {
	adminToken := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// Create a product and place an order to ensure at least one order exists
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "OrderGetProduct",
		"price":       25.0,
		"inStock":     8,
		"description": "desc",
		"category":    "OrderGetCat",
	}
	createReq.Var("input", input)
	tests.AuthRequest(createReq, adminToken)
	var createResp struct {
		CreateProduct struct{ ID string }
	}
	err := client.Run(context.TODO(), createReq, &createResp)
	require.NoError(t, err)
	productID := createResp.CreateProduct.ID

	customerToken := tests.Login(t, tests.CustomerEmail, tests.CustomerPassword)
	orderReq := graphql.NewRequest(`mutation($ids: [ID!]!) { placeOrder(productIds: $ids) { id } }`)
	orderReq.Var("ids", []string{productID})
	tests.AuthRequest(orderReq, customerToken)
	var orderResp struct {
		PlaceOrder struct{ ID string }
	}
	err = client.Run(context.TODO(), orderReq, &orderResp)
	require.NoError(t, err)
	orderID := orderResp.PlaceOrder.ID

	// Now, get the order
	getReq := graphql.NewRequest(`query($id: ID!) { order(id: $id) { id total status products { id name } user { id name } } }`)
	getReq.Var("id", orderID)
	tests.AuthRequest(getReq, customerToken)
	var getResp struct {
		Order struct {
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
	err = client.Run(context.TODO(), getReq, &getResp)
	require.NoError(t, err)
	require.Equal(t, orderID, getResp.Order.ID)
}
