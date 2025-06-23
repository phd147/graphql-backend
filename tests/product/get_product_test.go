package product

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestGetProduct(t *testing.T) {
	token := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// First, create a product to ensure it exists
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id name } }`)
	input := map[string]interface{}{
		"name":        "Test Product",
		"price":       10.5,
		"inStock":     5,
		"description": "desc",
		"category":    "TestCat",
	}
	createReq.Var("input", input)
	tests.AuthRequest(createReq, token)
	var createResp struct {
		CreateProduct struct {
			ID   string
			Name string
		}
	}
	err := client.Run(context.TODO(), createReq, &createResp)
	require.NoError(t, err)
	productID := createResp.CreateProduct.ID

	// Now, get the product
	getReq := graphql.NewRequest(`query($id: ID!) { product(id: $id) { id name price inStock description category } }`)
	getReq.Var("id", productID)
	tests.AuthRequest(getReq, token)
	var getResp struct {
		Product struct {
			ID          string
			Name        string
			Price       float64
			InStock     int
			Description string
			Category    string
		}
	}
	err = client.Run(context.TODO(), getReq, &getResp)
	require.NoError(t, err)
	require.Equal(t, productID, getResp.Product.ID)
}
