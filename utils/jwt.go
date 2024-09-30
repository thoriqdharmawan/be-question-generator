package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thoriqdharmawan/be-question-generator/api/v1/models/entity"
	"github.com/thoriqdharmawan/be-question-generator/config"
)

type Claims struct {
	User entity.User
	jwt.RegisteredClaims
}

func GenerateJWTToken(user entity.User) (string, error) {
	confVars, _ := config.New()
	var jwtSecretKey = []byte(confVars.JwtSecret)

	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWTToken(authToken string) (*Claims, error) {
	tokenString, errGetToken := GetTokenString(authToken)

	if errGetToken != nil {
		return nil, fmt.Errorf(errGetToken.Error())
	}

	confVars, _ := config.New()
	var jwtSecretKey = []byte(confVars.JwtSecret)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GetTokenString(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return token[7:], nil
}
