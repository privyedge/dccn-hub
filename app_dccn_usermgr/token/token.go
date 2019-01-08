package token

import (
	"errors"
	"time"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"

	"github.com/dgrijalva/jwt-go"
)

type TokenService interface {
	New(user *pb.User) (string, error)
	Verify(tokenString string) error
}

type Config struct {
	Issuer string `json:"issuer,omitempty"`
	Audience string `json:"audience,omitempty"`
	Subject string `json:"subject,omitempty"`
	ActiveTime int `json:"active_time,omitempty"`
	NotBefore int64 `json:"not_before,omitempty"`
	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	Secret string `json:"secret,omitempty"`
}

type Token struct {
	config *Config
}

// UserPayload is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type UserPayload struct {
	*pb.User
	jwt.StandardClaims
}

// New returns Token instance.
func New(conf *Config) *Token {
	return &Token{conf}
}

// New returns JWT string.
func (p *Token) New(user *pb.User) (string, error) {

	expireTime := time.Now().Add(time.Hour * time.Duration(p.config.ActiveTime)).Unix()

	// Create the Claims
	payload := UserPayload{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    p.config.Issuer,
			Subject:   p.config.Subject,
			NotBefore: p.config.NotBefore,
			Audience:  p.config.Audience,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(p.config.Secret)
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
	payload, ok := token.Claims.(*UserPayload)
	if ok && payload.User.Id !=  nil && token.Valid {
		return nil
	}
	return errors.New("invalid user")
}

