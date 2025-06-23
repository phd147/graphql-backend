package product

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestCreateProduct(t *testing.T) {
	token := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	client := tests.NewGraphQLClient()
	req := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id name price inStock description category } }`)
	input := map[string]interface{}{
		"name":        "Integration Product",
		"price":       99.99,
		"inStock":     10,
		"description": "Integration test product",
		"category":    "Integration",
	}
	req.Var("input", input)

	tests.AuthRequest(req, token)
	var resp struct {
		CreateProduct struct {
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
	require.Equal(t, input["name"], resp.CreateProduct.Name)
}

func TestCustomerCannotCreateProduct(t *testing.T) {
	token := tests.Login(t, tests.CustomerEmail, tests.CustomerPassword)
	client := tests.NewGraphQLClient()
	req := graphql.NewRequest(`mutation($input: CreateProductInput!) { createProduct(input: $input) { id } }`)
	input := map[string]interface{}{
		"name":        "ShouldFail",
		"price":       1.0,
		"inStock":     1,
		"description": "desc",
		"category":    "Cat",
	}
	req.Var("input", input)
	tests.AuthRequest(req, token)
	var resp struct {
		CreateProduct struct{ ID string }
	}
	err := client.Run(context.TODO(), req, &resp)
	require.Error(t, err)
}
