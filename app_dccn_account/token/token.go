package token

import (
	"github.com/Ankr-network/refactor/app_dccn_account/db_account"
	"github.com/Ankr-network/refactor/app_dccn_account/proto"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (

	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	key = []byte("mySuperSecretKeyLol")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	*accountmgr.Account
	jwt.StandardClaims
}

type TokenService interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *accountmgr.Account) (string, error)
}

type Token struct {
	dbaccount.AccountDBService
}

// Decode a token string into a token object
func (srv *Token) Decode(tokenString string) (*CustomClaims, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode a claim into a JWT
func (srv *Token) Encode(ac *accountmgr.Account) (string, error) {

	expireToken := time.Now().Add(time.Hour * 72).Unix()

	// Create the Claims
	claims := CustomClaims{
		ac,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "ankr",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(key)
}

