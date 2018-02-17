package authentication

import (
	// standard libraries
	"fmt"
	"time"
	// custom handlers
	user "github.com/callistom/go-graphql-auth/structs"
	jwt "github.com/dgrijalva/jwt-go"
)

// MyCustomClaims jwt claims with id
type MyCustomClaims struct {
	ID uint
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
