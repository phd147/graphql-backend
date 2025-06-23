package product

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestGetProducts(t *testing.T) {
	token := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	// Create a product to ensure at least one exists
	adminToken := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	createReq := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "ListProduct",
		"price":       5.0,
		"inStock":     3,
		"description": "desc",
		"category":    "ListCat",
	}
	createReq.Var("input", input)
	tests.AuthRequest(createReq, adminToken)
	_ = client.Run(context.TODO(), createReq, &struct{ CreateProduct struct{ ID string } }{})

	req := graphql.NewRequest(`query { products(limit: 10, offset: 0) { id name price inStock description category } }`)
	tests.AuthRequest(req, token)
	var resp struct {
		Products []struct {
			ID          string
			Name        string
			Price       float64
			InStock     int
			Description string
			Category    string
		}
	}
	err := client.Run(context.TODO(), req, &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Products)
}
