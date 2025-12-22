package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret    string
	AccessTTL time.Duration
	Issuer    string
}

type Claims struct {
	UserID uint     `json:"uid"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func SignAccessToken(cfg JWTConfig, userID uint, roles []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTTL)),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(cfg.Secret))
}

func VerifyAccessToken(cfg JWTConfig, tokenStr string) (*Claims, error) {
	keyFn := func(t *jwt.Token) (any, error) {
		// Asegurar algoritmo
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return []byte(cfg.Secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, keyFn, jwt.WithIssuer(cfg.Issuer))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}
