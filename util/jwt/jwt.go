package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)



type MyCustomClaims struct {
	Token string `json:"token"`
	UserGUID string `json:"UserGUID"`
	TokenExpiry int64 `json:"TokenExpiry"`
	jwt.StandardClaims
}


var secret = "28iOiJiYXIiLCJleHAiOjE1MDAwLCJ"

func CreateJwtToken(value string, name string) string {
	mySigningKey := []byte(secret)

	// Create the Claims
	claims := MyCustomClaims{
		value,
		name,
		86400,
		jwt.StandardClaims{
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	//fmt.Printf("%v %v", ss, err)
	return ss
}

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


