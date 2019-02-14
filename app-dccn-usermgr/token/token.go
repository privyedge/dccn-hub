package token

import (
	"errors"
	"log"
	"time"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	jwt "github.com/dgrijalva/jwt-go"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
)

var secret = []byte(ankr_default.Secret)

type IToken interface {
	NewAuthenticationToken(user *usermgr.User) (string, error)
	NewAuthorizationToken(email string) (string, error)
	VerifyAuthenticationToken(tokenString string) (*AuthenticationPayload, error)
	VerifyAuthorizationToken(tokenString string) (*AuthorizationPayload, error)
	VerifyAndRefreshAuthenticationToken(tokenString string) (string, error)
}

type Token struct {
	RefreshTokenValidTime  int
	AccessTokenValidTime   int
	ActivateUserValidTime  int
	ResetPasswordValidTime int
}

// AuthenticationPayload is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type AuthenticationPayload struct {
	user *usermgr.User
	jwt.StandardClaims
}

// AuthorizationPayload used to reset password
type AuthorizationPayload struct {
	Email string
	jwt.StandardClaims
}

// New returns Token instance.
func New() *Token {
	return &Token{
		AccessTokenValidTime:   ankr_default.AccessTokenValidTime,
		RefreshTokenValidTime:  ankr_default.RefreshTokenValidTime,
		ActivateUserValidTime:  ankr_default.ActivateCodeValidTime,
		ResetPasswordValidTime: ankr_default.ActivateCodeValidTime,
	}
}

// New returns JWT string.
func (p *Token) NewAuthenticationToken(user *usermgr.User) (string, error) {

	expireTime := time.Now().Add(time.Minute * time.Duration(p.AccessTokenValidTime)).Unix()

	// Create the Claims
	payload := AuthenticationPayload{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "ankr.network",
			Id:        user.Id,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(secret)
}

func (p *Token) NewAuthorizationToken(email string) (string, error) {

	expireTime := time.Now().Add(time.Minute * time.Duration(p.ResetPasswordValidTime)).Unix()

	// Create the Claims
	payload := AuthorizationPayload{
		email,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "ankr.network",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(secret)
}

// VerifyAuthentication validates authentication token
func (p *Token) VerifyAuthenticationToken(tokenString string) (*AuthenticationPayload, error) {

	log.Println("Debug into Verify: ", tokenString)
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &AuthenticationPayload{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Validate the token
	if payload, ok := token.Claims.(*AuthenticationPayload); ok && token.Valid {
		return payload, nil
	}
	return nil, errors.New("token is invalid")
}

// VerifyAuthorizationToken validates authorization token
func (p *Token) VerifyAuthorizationToken(tokenString string) (*AuthorizationPayload, error) {

	log.Println("Debug into Verify: ", tokenString)
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &AuthorizationPayload{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// Validate the token
	if payload, ok := token.Claims.(*AuthorizationPayload); ok && token.Valid {
		return payload, nil
	}
	return nil, errors.New("token is invalid")
}

func (p *Token) VerifyAndRefreshAuthenticationToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	payload, err := p.VerifyAuthenticationToken(tokenString)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	jwt.TimeFunc = time.Now
	payload.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(p.AccessTokenValidTime) * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// Sign token and return
	return token.SignedString(secret)
}
