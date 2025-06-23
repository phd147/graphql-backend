package user

import (
	"context"
	"github.com/machinebox/graphql"
	"testing"

	"github.com/stretchr/testify/require"

	"graphql-backend/tests"
)

func TestGetMe(t *testing.T) {
	token := tests.Login(t, tests.CustomerEmail, tests.CustomerPassword)
	client := tests.NewGraphQLClient()
	req := graphql.NewRequest(`query { me { id name email } }`)
	tests.AuthRequest(req, token)
	var resp struct {
		Me struct {
			ID    string
			Name  string
			Email string
		}
	}
	err := client.Run(context.TODO(), req, &resp)
	require.NoError(t, err)
	require.Equal(t, tests.CustomerEmail, resp.Me.Email)
}
