package security

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"plexus-bff-service-go/internal/app/config"
)

func TestValidatorFallsBackToUnverifiedWhenNoJWKConfigured(t *testing.T) {
	validator := NewValidator(config.SecurityJWTConfig{
		AuthoritiesClaim: "cognito:groups",
		RolePrefix:       "ROLE_",
	})

	token := fakeUnsignedJWT(t, map[string]any{
		"cognito:groups": []string{"ADMIN"},
		"scope":          "admin profile",
	})

	auth, err := validator.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if !HasAuthority(auth, "ROLE_ADMIN") || !HasAuthority(auth, "SCOPE_admin") {
		t.Fatalf("expected authorities to be extracted")
	}
}

func TestValidatorValidatesAgainstJWKSetAndIssuer(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]any{
				{
					"kty": "RSA",
					"kid": "test-key",
					"use": "sig",
					"alg": "RS256",
					"n":   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
				},
			},
		})
	}))
	defer jwks.Close()

	validator := NewValidatorWithClient(config.SecurityJWTConfig{
		JWKSetURI:       jwks.URL,
		IssuerURI:       "https://issuer.example.com",
		AuthoritiesClaim: "cognito:groups",
		RolePrefix:      "ROLE_",
	}, jwks.Client())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":            "https://issuer.example.com",
		"exp":            time.Now().Add(time.Hour).Unix(),
		"cognito:groups": []string{"ADMIN"},
	})
	token.Header["kid"] = "test-key"
	signed, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	auth, err := validator.Parse(signed)
	if err != nil {
		t.Fatalf("validate token: %v", err)
	}
	if !HasAuthority(auth, "ROLE_ADMIN") {
		t.Fatalf("expected validated token to contain ROLE_ADMIN")
	}
}

func TestValidatorRejectsWrongIssuer(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}

	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]any{
				{
					"kty": "RSA",
					"kid": "test-key",
					"use": "sig",
					"alg": "RS256",
					"n":   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
				},
			},
		})
	}))
	defer jwks.Close()

	validator := NewValidatorWithClient(config.SecurityJWTConfig{
		JWKSetURI: jwks.URL,
		IssuerURI: "https://issuer.example.com",
	}, jwks.Client())

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://other.example.com",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = "test-key"
	signed, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	if _, err := validator.Parse(signed); err == nil {
		t.Fatalf("expected issuer validation failure")
	}
}

func fakeUnsignedJWT(t *testing.T, claims map[string]any) string {
	t.Helper()
	header, err := json.Marshal(map[string]any{"alg": "none", "typ": "JWT"})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	body, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(body) + "."
}
