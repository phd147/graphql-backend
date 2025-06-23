package http_transport

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt/v5"
	"graphql-backend/graph/model"
	"net/http"
	"strings"
)

type contextKey string

const (
	UserContextKey contextKey = "user"
)

// UserClaims represents the JWT claims
type UserClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	*jwt.RegisteredClaims
}

// AuthMiddleware is a middleware for authentication
func AuthMiddleware(jwtHandler JwtHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if the header has the Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				next.ServeHTTP(w, r)
				return
			}

			// Parse the token
			token, err := jwt.ParseWithClaims(parts[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
				kid, ok := token.Header["kid"]
				if !ok {
					return nil, errors.New("missing kid in token header")
				}

				// Verify the token's signing method
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				// Get the public key from the JWT handler
				publicKey, err := jwtHandler.GetPublicKey(r.Context(), kid.(string))
				if err != nil {
					return nil, errors.New("invalid token")
				}

				return publicKey, nil
			})

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
				// Add the user to the context
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext gets the user from the context
func GetUserFromContext(ctx context.Context) *UserClaims {
	user, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok {
		return nil
	}
	return user
}

var HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
	user := GetUserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	if user.Role != string(role) {
		return nil, fmt.Errorf("user does not have the required role: %s", role)
	}

	// or let it pass through
	return next(ctx)
}

var HasAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	user := GetUserFromContext(ctx)
	if user == nil {
		return nil, errors.New("unauthorized")
	}

	// or let it pass through
	return next(ctx)
}
