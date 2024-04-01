package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

// DecodeJWT - декодирует JWT токен в структуру
func DecodeJWT(token string, jwtKey []byte, model jwt.Claims) error {
	token = strings.ReplaceAll(token, "Bearer ", "")

	tkn, err := jwt.ParseWithClaims(token, model, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return err
		}

		return err
	}

	if !tkn.Valid {
		return err
	}

	return nil
}
