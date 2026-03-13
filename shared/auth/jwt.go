package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     int64  `json:"user_id"`
	EmployeeID string `json:"employee_id"`
	JTI        string `json:"jti"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(secret string, userID int64, employeeID, jti string, ttl time.Duration) (string, int64, error) {
	expireAt := time.Now().Add(ttl)

	claims := Claims{
		UserID:     userID,
		EmployeeID: employeeID,
		JTI:        jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   employeeID,
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signed, expireAt.Unix(), nil
}

func ParseAccessToken(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
