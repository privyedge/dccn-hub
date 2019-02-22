package user_util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/micro/go-micro/metadata"
)

type Token struct {
	Exp int64
	Jti string
	Iss string
}

func GetUserID(ctx context.Context) string {
	meta, ok := metadata.FromContext(ctx)
	// Note this is now uppercase (not entirely sure why this is...)
	var token string
	if ok {
		token = meta["token"]
	}
	parts := strings.Split(token, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("decode error:", err)

	}
	fmt.Println(string(decoded))
	var dat Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		fmt.Println("Unmarshal error:", err)
	}
	return string(dat.Jti)
}

func VerifyAccessToken(refreshToken string) (string, error) {
	parts := strings.Split(refreshToken, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	now := time.Now().Unix()

	if now > int64(dat.Exp) {
		return "", ankr_default.ErrTokenParseFailed
	}

	return string(dat.Jti), nil
}

func GetIdFromToken(refreshToken string) (string, error) {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 3 {
		return "", ankr_default.ErrTokenParseFailed
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	return string(dat.Jti), nil
}
