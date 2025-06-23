package user

import (
	"context"
	"github.com/machinebox/graphql"
	"testing"

	"github.com/stretchr/testify/require"
	"graphql-backend/tests"
)

func TestLogin_Success(t *testing.T) {
	token := tests.Login(t, tests.AdminEmail, tests.AdminPassword)
	require.NotEmpty(t, token)
}

func TestLogin_Fail(t *testing.T) {
	client := tests.NewGraphQLClient()
	request := graphql.NewRequest(`mutation($input: LoginInput!) {  login(input: $input) {    accessToken  }}`)
	request.Var("input", map[string]interface{}{
		"email":    "wrong@example.com",
		"password": "wrong",
	})
	var resp struct {
		Login struct {
			AccessToken string `json:"accessToken"`
		}
	}
	err := client.Run(context.TODO(), request, &resp)
	require.Error(t, err)
}
