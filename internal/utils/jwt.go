package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key")

type Claims struct {
	UserID int
	Email  string
	Role   string
	jwt.RegisteredClaims
}

func GetTokenWithoutBearer(authToken string) string {
	if len(authToken) > 7 && authToken[:7] == "Bearer " {
		return authToken[7:]
	}
	return authToken
}

func GenerateJWT(userID int, email string, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func IsValidJWT(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false, err
	}

	return true, nil
}

func GetClaimsFromJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

func GetClaimsFromJWTUsingContext(c *gin.Context) (*Claims, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return nil, jwt.ErrSignatureInvalid
	}

	token = GetTokenWithoutBearer(token)
	claims, err := GetClaimsFromJWT(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
