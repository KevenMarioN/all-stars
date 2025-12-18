package middlewares_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/KevenMarioN/all-stars/server/middlewares"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	ID   int
	Role string
}

func TestAuthMiddleware(t *testing.T) {
	exp := func() time.Time { return time.Now().Add(24 * time.Hour) }
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	publicKey := &privateKey.PublicKey
	// TODO: Need implementation others scenarios!
	scenarios := []struct {
		name                 string
		data                 Auth
		auth                 *middlewares.AuthMiddleware[Auth]
		expectedInvalidToken bool
	}{
		{
			name: "Valid Token HMAC",
			data: Auth{
				ID:   1,
				Role: "admin",
			},
			auth: middlewares.NewAuthMiddleware[Auth](jwt.GetSigningMethod(jwt.SigningMethodHS384.Name), []byte("secret"), []byte("secret"), exp, nil),
		},
		{
			name: "Invalid Token HMAC",
			data: Auth{
				ID:   30,
				Role: "user",
			},
			auth:                 middlewares.NewAuthMiddleware[Auth](jwt.GetSigningMethod(jwt.SigningMethodHS384.Name), []byte("secret"), []byte("terces"), exp, nil),
			expectedInvalidToken: true,
		},
		{
			name: "Valid Token ECDSA",
			data: Auth{
				ID:   2,
				Role: "user",
			},
			auth: middlewares.NewAuthMiddleware[Auth](jwt.SigningMethodES256, privateKey, publicKey, exp, nil),
		},
	}

	for _, tt := range scenarios {
		t.Run(tt.name, func(t *testing.T) {
			token, err := tt.auth.CreateToken(tt.data)
			if err != nil {
				t.Errorf("CreateToken() error = %v", err)
			}
			claims := &middlewares.AuthClaims[Auth]{}
			if err = tt.auth.ParseToken(token, claims); err != nil {
				if !tt.expectedInvalidToken {
					t.Errorf("ParseToken() error = %v", err)
				}
			}
			if claims.Payload != tt.data {
				t.Errorf("Payload = %v, want %v", claims.Payload, tt.data)
			}
			if claims.ExpiresAt.Before(time.Now()) {
				t.Errorf("ExpiresAt = %v, want after %v", claims.ExpiresAt, time.Now())
			}
			if claims.IssuedAt.Before(time.Now().Add(-time.Hour)) {
				t.Errorf("IssuedAt = %v, want after %v", claims.IssuedAt, time.Now().Add(-time.Hour))
			}
		})
	}
}
