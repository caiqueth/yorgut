package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type Token struct {
	iss            string
	azp            string
	aud            string
	sub            string
	email          string
	email_verified bool
	nbf            int
	name           string
	picture        string
	given_name     string
	family_name    string
	iat            int
	exp            int
	jti            string
}

func parseToken(token string) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	payload := parts[1]

	if mod := len(payload) % 4; mod != 0 {
		payload += strings.Repeat("=", 4-mod)
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("invalid base64 encoding")
	}

	var claims *Token
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, errors.New("invalid JSON format")
	}

	return claims, nil
}

func GetJWTToken(c *gin.Context) (*Token, bool) {
	token, err := parseToken(c.GetString("user_token"))
	if err != nil {
		return nil, false
	}
	return token, true
}

func SetJWTToken(c *gin.Context, token string) {
	c.Set("user_token", token)
}

func GetEmailFromToken(c *gin.Context) (string, bool) {
	token, ok := GetJWTToken(c)
	if !ok {
		return "", false
	}
	return token.email, true
}
