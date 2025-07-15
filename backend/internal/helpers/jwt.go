package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"yorgut/internal/model"

	"github.com/gin-gonic/gin"
)

var USER_TOKEN_KEY string = "user_token"

func ParseToken(token string) (*model.Token, error) {
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

	var claims *model.Token
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, errors.New("invalid JSON format")
	}

	return claims, nil
}

func SetJWTToken(c *gin.Context, token map[string]interface{}) {
	c.Set(USER_TOKEN_KEY, token)
}

func GetTokenFromContext(c *gin.Context) (*model.Token, bool) {
	data, exists := c.Get(USER_TOKEN_KEY)
	if !exists {
		return nil, false
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var token model.Token
	if err := json.Unmarshal(jsonBytes, &token); err != nil {
		return nil, false
	}

	return &token, true
}

func GetUserEmail(c *gin.Context) string {
	tkn, ok := GetTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "could not find auth token"})
		return ""
	}

	return tkn.Email
}
