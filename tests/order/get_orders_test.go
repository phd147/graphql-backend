package order

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestGetOrders(t *testing.T) {
	adminToken := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// Create a product and place an order to ensure at least one order exists
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "OrderListProduct",
		"price":       15.0,
		"inStock":     7,
		"description": "desc",
		"category":    "OrderListCat",
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
	_ = client.Run(context.TODO(), orderReq, &struct{ PlaceOrder struct{ ID string } }{})

	// Now, get orders
	getReq := graphql.NewRequest(`query { orders(limit: 10, offset: 0) { id total status products { id name } user { id name } } }`)
	tests.AuthRequest(getReq, customerToken)
	var getResp struct {
		Orders []struct {
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
	require.NotEmpty(t, getResp.Orders)
}
