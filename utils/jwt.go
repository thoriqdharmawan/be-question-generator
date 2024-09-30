package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
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

func VerifyJWTToken(tokenString string) (*Claims, error) {
	godotenv.Load(".env")
	var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

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

func VerifyJWTTokenHandler(token string) error {
	if token == "" {
		return fmt.Errorf("missing authorization header")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return fmt.Errorf("invalid authorization header format")
	}

	tokenString := token[7:]

	_, err := VerifyJWTToken(tokenString)

	if err != nil {
		return fmt.Errorf("invalid token")
	}

	return nil
}
