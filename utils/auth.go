package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(email string, seconds int) (string, error) {
	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	if seconds != 0 {
		expDuration := time.Duration(seconds) * time.Second
		exp := now.Add(expDuration).Unix()
		claims["exp"] = exp
	}
	claims["email"] = email
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenFromHeader(authorizationHeader string) string {
	token, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
	if ok {
		return token
	}
	return ""
}

func ValidateSigningMethod(jwtToken *jwt.Token) error {
	if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
		return fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
	}
	return nil
}

func GetSecretKeyAPIKey(apiKey string) ([]byte, error) {
	return []byte(apiKey), nil
}

func ParseAndValidateToken(tokenString, apiKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if err := ValidateSigningMethod(jwtToken); err != nil {
			return nil, err
		}
		return GetSecretKeyAPIKey(apiKey)
	})
}
