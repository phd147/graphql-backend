package product

import (
	"context"
	"github.com/machinebox/graphql"
	"testing"

	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestUpdateProduct(t *testing.T) {
	token := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// Create a product to update
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "ToUpdate",
		"price":       1.0,
		"inStock":     1,
		"description": "desc",
		"category":    "Cat",
	}
	createReq.Var("input", input)
	tests.AuthRequest(createReq, token)
	var createResp struct {
		CreateProduct struct{ ID string }
	}
	err := client.Run(context.TODO(), createReq, &createResp)
	require.NoError(t, err)
	productID := createResp.CreateProduct.ID

	// Update the product
	updateReq := graphql.NewRequest(`mutation($input: UpdateProductInput!) { updateProduct(input: $input) { id name price inStock description category } }`)
	updateInput := map[string]interface{}{
		"id":          productID,
		"name":        "UpdatedName",
		"price":       2.0,
		"inStock":     2,
		"description": "updated desc",
		"category":    "UpdatedCat",
	}
	updateReq.Var("input", updateInput)
	tests.AuthRequest(updateReq, token)
	var updateResp struct {
		UpdateProduct struct {
			ID          string
			Name        string
			Price       float64
			InStock     int
			Description string
			Category    string
		}
	}
	err = client.Run(context.TODO(), updateReq, &updateResp)
	require.NoError(t, err)
	require.Equal(t, updateInput["name"], updateResp.UpdateProduct.Name)
}

func TestCustomerCannotUpdateProduct(t *testing.T) {
	token := tests.Login(t, tests.CustomerEmail, tests.CustomerPassword)
	client := tests.NewGraphQLClient()
	// Try to update a (likely non-existent) product
	updateReq := graphql.NewRequest(`mutation($input: UpdateProductInput!) { updateProduct(input: $input) { id } }`)
	updateInput := map[string]interface{}{
		"id":   "non-existent-id",
		"name": "ShouldFail",
	}
	updateReq.Var("input", updateInput)
	tests.AuthRequest(updateReq, token)
	var updateResp struct {
		UpdateProduct struct{ ID string }
	}
	err := client.Run(context.TODO(), updateReq, &updateResp)
	require.Error(t, err)
}
