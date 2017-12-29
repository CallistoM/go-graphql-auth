package authentication

import (
	"fmt"
	"time"

	user "github.com/callistom/jwt_graphql_server/structs"
	jwt "github.com/dgrijalva/jwt-go"
)

type MyCustomClaims struct {
	ID string
	jwt.MapClaims
}

//PublicKey set secret public key
var PublicKey = []byte("secret")

// GenerateToken generates JWT token en returns it
func GenerateToken(user user.User) (string, error) {

	// Create the Claims
	claims := MyCustomClaims{
		user.ID,
		jwt.MapClaims{
			"exp": time.Now().Add(time.Hour * time.Duration(24)).Unix(),
			"iat": time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(PublicKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// CheckToken checks if token is valid else returns error
func CheckToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return PublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		return token, nil
	}

	return token, nil
}
