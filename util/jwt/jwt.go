package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)



type MyCustomClaims struct {
	Token string `json:"token"`
	UserGUID string `json:"UserGUID"`
	TokenExpiry int64 `json:"TokenExpiry"`
	jwt.StandardClaims
}
type JWTUserTokenInfo struct {
	UserGUID    string   `json:"user_id"`
	UserName    string   `json:"user_name"`
	TokenExpiry int64    `json:"exp"`
	Scope       []string `json:"scope"`
	jwt.StandardClaims
}


var secret = "28iOiJiYXIiLCJleHAiOjE1MDAwLCJ"

func CreateJwtToken(value string, name string) string {
	mySigningKey := []byte(secret)

	// Create the Claims
	claims := JWTUserTokenInfo{
		value,
		name,
		time.Now().Unix() + 86400,
		[]string{"stratos.admin"},
		jwt.StandardClaims{},
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

	if claims, ok := token.Claims.(*JWTUserTokenInfo); ok && token.Valid {
		//fmt.Printf("%v %v", claims.Token, claims.StandardClaims.ExpiresAt)
		return claims.UserGUID
	} else {
		fmt.Println(err)
		return ""
	}


}


