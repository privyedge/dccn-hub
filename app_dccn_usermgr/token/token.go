package token

import (
	"errors"
	"time"

	"github.com/Ankr-network/refactor/proto/usermgr"
	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	New(user *usermgr.User) (string, error)
	Verify(tokenString string) error
}

type Config struct {
	Issuer string
	Audience string
	Subject string
	ActiveTime int
	NotBefore int64
	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	Secret string
}

type Token struct {
	config *Config
}

// UserPayload is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type UserPayload struct {
	*usermgr.User
	jwt.StandardClaims
}

func New(conf *Config) *Token {
	return &Token{conf}
}

// Verify a token string into a token object
func (p *Token) Verify(tokenString string) error {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (interface{}, error) {
		return p.config.Secret, nil
	})

	if err != nil {
		return err
	}

	// Validate the token
	if _, ok := token.Claims.(*UserPayload); ok && token.Valid {
		return nil
	} else {
		return errors.New("invalid taskmgr")
	}
}

// New a claim into a JWT
func (p *Token) New(user *usermgr.User) (string, error) {

	expireToken := time.Now().Add(time.Minute * time.Duration(p.config.ActiveTime)).Unix()

	// Create the Claims
	payload := UserPayload{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    p.config.Issuer,
			Subject:p.config.Subject,
			NotBefore:p.config.NotBefore,
			Audience:p.config.Audience,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(p.config.Secret)
}

