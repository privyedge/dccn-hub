package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)



type MyCustomClaims struct {
	Token string `json:"foo"`
	jwt.StandardClaims
}


var secret = "28iOiJiYXIiLCJleHAiOjE1MDAwLCJ"

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func CreateJwtToken(value string) string {
	mySigningKey := []byte(secret)

	// Create the Claims
	claims := MyCustomClaims{
		value,
		jwt.StandardClaims{
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	//fmt.Printf("%v %v", ss, err)
	return ss
}

// Example creating a token using a custom claims type.  The StandardClaim is embedded
// in the custom type to allow for easy encoding, parsing and validation of standard claims.
func ParseJwtToken(tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Token, claims.StandardClaims.ExpiresAt)
		return claims.Token
	} else {
		fmt.Println(err)
		return ""
	}


}


