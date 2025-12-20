package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const AUTH_KEY_PAYLOAD MwKey = "payload"

var ErrInvalidToken = errors.New("invalid token")

type AuthClaims[T any] struct {
	Payload T `json:"payload"`
	jwt.RegisteredClaims
}

type AuthMiddleware[T any] struct {
	ExpiresAt time.Time
	NotBefore *time.Time
	signKey   any
	verifyKey any
	method    jwt.SigningMethod
}

func NewAuthMiddleware[T any](method jwt.SigningMethod, signKey any, verifyKey any, expiresAt func() time.Time, notBefore func() time.Time) *AuthMiddleware[T] {
	var (
		auth = &AuthMiddleware[T]{
			signKey:   signKey,
			verifyKey: verifyKey,
			method:    method,
		}
		notBeforeTime time.Time
	)
	if expiresAt == nil {
		auth.ExpiresAt = time.Now().Add(24 * time.Hour)
	} else {
		auth.ExpiresAt = expiresAt()
	}

	if notBefore != nil {
		notBeforeTime = notBefore()
		auth.NotBefore = &notBeforeTime
	}

	return auth
}

func (m *AuthMiddleware[T]) CreateToken(data T) (string, error) {
	claims := AuthClaims[T]{
		Payload: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(m.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	if m.NotBefore != nil {
		claims.NotBefore = jwt.NewNumericDate(*m.NotBefore)
	}
	token := jwt.NewWithClaims(m.method, claims)
	return token.SignedString(m.signKey)
}

func (m *AuthMiddleware[T]) ParseToken(tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != m.method.Alg() {
			return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}
		return m.verifyKey, nil
	})

	if err != nil || !token.Valid {
		return ErrInvalidToken
	}

	return nil
}

func (m *AuthMiddleware[T]) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		var claims AuthClaims[T]
		if err := m.ParseToken(tokenString, &claims); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AUTH_KEY_PAYLOAD, claims.Payload)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
