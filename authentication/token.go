package authentication

import (
	"fmt"
	"time"

	user "github.com/callistom/api-project/structs"
	jwt "github.com/dgrijalva/jwt-go"
)

//PublicKey set secret public key
var PublicKey = []byte("secret")

func GenerateToken(user user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(24)).Unix(),
		"iat": time.Now().Unix(),
		"sub": user.ID,
	})
	tokenString, err := token.SignedString(PublicKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func CheckToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
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
