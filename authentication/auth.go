package authentication

import (
	// standard libraries
	"context"
	"fmt"
	"net/http"
)

type contextKey string

// read context key
func (c contextKey) String() string {
	return string(c)
}

// Auth checks if user is authenticated
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		token := r.Header.Get("Authorization")

		jwt, err := CheckToken(token)

		if err != nil {
			fmt.Println(err)
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, "jwt", jwt)))
	})
}
