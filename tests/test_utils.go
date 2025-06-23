package tests

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
	"testing"
)

var serverURL = "http://localhost:8080"

const (
	AdminEmail       = "admin@example.com"
	AdminPassword    = "secret"
	CustomerEmail    = "customer@example.com"
	CustomerPassword = "secret"
)

func NewGraphQLClient() *graphql.Client {
	return graphql.NewClient(fmt.Sprintf("%s/query", serverURL))
}

func Login(t *testing.T, email, password string) string {
	client := NewGraphQLClient()
	request := graphql.NewRequest(`mutation($input: LoginInput!) {  login(input: $input) {    accessToken  }}`)
	request.Var("input", map[string]interface{}{
		"email":    email,
		"password": password,
	})
	var resp struct {
		Login struct {
			AccessToken string `json:"accessToken"`
		}
	}
	err := client.Run(context.Background(), request, &resp)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	return resp.Login.AccessToken
}

func AuthRequest(req *graphql.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
}

func init() {
	if v := os.Getenv("SERVER_URL"); v != "" {
		serverURL = v
	}
}
