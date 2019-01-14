package token

import (
	"errors"
	"time"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/v1"

	jwt "github.com/dgrijalva/jwt-go"
)

var secret = []byte("14444749c1ecc982cd0f91113db98211")

type IToken interface {
	New(user *pb.User) (string, error)
	Verify(tokenString string) error
}

type Token struct {
	activeTime int
}

// UserPayload is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type UserPayload struct {
	user *pb.User
	jwt.StandardClaims
}

// New returns Token instance.
func New(activeTime int) *Token {
	return &Token{activeTime}
}

// New returns JWT string.
func (p *Token) New(user *pb.User) (string, error) {

	expireTime := time.Now().Add(time.Minute * time.Duration(p.activeTime)).Unix()

	// Create the Claims
	payload := UserPayload{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "ankr_network",
			// Subject:   p.config.Subject,
			// NotBefore: p.config.NotBefore,
			// Audience:  p.config.Audience,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(secret)
}

// Verify a token string into a token object
func (p *Token) Verify(tokenString string) error {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	// Validate the token
	_, ok := token.Claims.(*UserPayload)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid user")
}
