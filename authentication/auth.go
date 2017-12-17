package authentication

import (
	"context"
	"fmt"
	"net/http"
)

type contextKey string

// read context key
func (c contextKey) String() string {
	return string(c)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		jwt, err := CheckToken(token)
		if err != nil {
			fmt.Println(err)
		}

		// set context key
		var (
			contextKeyAuthtoken = contextKey("jwt")
		)

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, contextKeyAuthtoken, jwt)))
	})
}
