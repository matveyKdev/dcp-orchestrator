package access_token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

var secretKey = []byte("secretClusterKey")

func GenerateToken(secretKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string, clusterToken string) (bool, error) {
	tokenInput := strings.Trim(tokenString, " ")
	if tokenInput == "" {
		return false, fmt.Errorf("token cannot be empty")
	}

	if tokenInput != clusterToken {
		return false, nil
	}
	return true, nil
}
