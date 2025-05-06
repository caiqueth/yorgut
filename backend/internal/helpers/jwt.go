package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var USER_TOKEN_KEY string = "user_token"

type Token struct {
	Iss            string
	Azp            string
	Aud            string
	Sub            string
	Email          string
	Email_verified bool
	Nbf            int
	Name           string
	Picture        string
	Given_name     string
	Family_name    string
	Iat            int
	Exp            int
	Jti            string
}

func ParseToken(token string) (*Token, error) {
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

func ValidateToken(token Token) Token {
	return token
}

func SetJWTToken(c *gin.Context, token *Token) {
	c.Set(USER_TOKEN_KEY, token)
}

func GetTokenFromContext(c *gin.Context) (*Token, bool) {
	value, exists := c.Get(USER_TOKEN_KEY)
	if !exists {
		return nil, false
	}

	token, ok := value.(*Token)
	return token, ok
}

func GetUserEmail(c *gin.Context) string {
	tkn, ok := GetTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not find auth token"})
		return ""
	}

	return tkn.Email
}
