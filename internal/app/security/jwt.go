package security

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"

	"plexus-bff-service-go/internal/app/config"
)

type Validator struct {
	cfg        config.SecurityJWTConfig
	httpClient *http.Client
	mu         sync.RWMutex
	keySet     jwkSet
}

type jwkSet struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	KeyType string `json:"kty"`
	KeyID   string `json:"kid"`
	Use     string `json:"use"`
	Alg     string `json:"alg"`
	N       string `json:"n"`
	E       string `json:"e"`
}

func NewValidator(cfg config.SecurityJWTConfig) *Validator {
	return &Validator{
		cfg:        cfg,
		httpClient: http.DefaultClient,
	}
}

func NewValidatorWithClient(cfg config.SecurityJWTConfig, client *http.Client) *Validator {
	if client == nil {
		client = http.DefaultClient
	}
	return &Validator{
		cfg:        cfg,
		httpClient: client,
	}
}

func (v *Validator) Parse(token string) (*Authentication, error) {
	if v.cfg.JWKSetURI == "" && v.cfg.IssuerURI == "" {
		return parseUnverified(token, v.cfg)
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256", "RS384", "RS512"}))
	claims := jwt.MapClaims{}
	parsed, err := parser.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method == nil {
			return nil, errors.New("missing signing method")
		}
		keyID, _ := t.Header["kid"].(string)
		if keyID == "" {
			return nil, errors.New("missing kid header")
		}
		return v.lookupKey(context.Background(), keyID)
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid jwt")
	}
	if err := validateIssuer(claims, v.cfg.IssuerURI); err != nil {
		return nil, err
	}
	return newAuthentication(token, claims, v.cfg), nil
}

func parseUnverified(token string, cfg config.SecurityJWTConfig) (*Authentication, error) {
	parser := jwt.NewParser()
	claims := jwt.MapClaims{}
	if _, _, err := parser.ParseUnverified(token, claims); err != nil {
		return nil, err
	}
	return newAuthentication(token, claims, cfg), nil
}

func newAuthentication(token string, claims jwt.MapClaims, cfg config.SecurityJWTConfig) *Authentication {
	return &Authentication{
		Claims:      map[string]any(claims),
		Authorities: extractAuthorities(claims, cfg),
		Token:       token,
	}
}

func extractAuthorities(claims jwt.MapClaims, cfg config.SecurityJWTConfig) []string {
	authorities := make([]string, 0)
	claimName := cfg.AuthoritiesClaim
	if claimName == "" {
		claimName = "cognito:groups"
	}
	rolePrefix := cfg.RolePrefix
	if rolePrefix == "" {
		rolePrefix = "ROLE_"
	}

	switch groups := claims[claimName].(type) {
	case []any:
		for _, group := range groups {
			if text, ok := group.(string); ok && text != "" {
				authorities = append(authorities, rolePrefix+text)
			}
		}
	case []string:
		for _, group := range groups {
			if group != "" {
				authorities = append(authorities, rolePrefix+group)
			}
		}
	}

	switch scopes := claims["scope"].(type) {
	case string:
		for _, scope := range strings.Split(scopes, " ") {
			scope = strings.TrimSpace(scope)
			if scope != "" {
				authorities = append(authorities, "SCOPE_"+scope)
			}
		}
	case []any:
		for _, scope := range scopes {
			if text, ok := scope.(string); ok && text != "" {
				authorities = append(authorities, "SCOPE_"+text)
			}
		}
	}
	return authorities
}

func validateIssuer(claims jwt.MapClaims, expected string) error {
	if expected == "" {
		return nil
	}
	issuer, _ := claims["iss"].(string)
	if issuer == "" {
		return errors.New("missing issuer claim")
	}
	if issuer != expected {
		return fmt.Errorf("unexpected issuer %q", issuer)
	}
	return nil
}

func (v *Validator) lookupKey(ctx context.Context, keyID string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	for _, key := range v.keySet.Keys {
		if key.KeyID == keyID {
			v.mu.RUnlock()
			return toRSAPublicKey(key)
		}
	}
	v.mu.RUnlock()

	if err := v.refreshKeys(ctx); err != nil {
		return nil, err
	}

	v.mu.RLock()
	defer v.mu.RUnlock()
	for _, key := range v.keySet.Keys {
		if key.KeyID == keyID {
			return toRSAPublicKey(key)
		}
	}
	return nil, fmt.Errorf("kid %q not found in jwk set", keyID)
}

func (v *Validator) refreshKeys(ctx context.Context) error {
	if v.cfg.JWKSetURI == "" {
		return errors.New("jwk-set-uri is not configured")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.cfg.JWKSetURI, nil)
	if err != nil {
		return err
	}
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("jwk fetch returned %d", resp.StatusCode)
	}
	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return err
	}
	v.mu.Lock()
	v.keySet = set
	v.mu.Unlock()
	return nil
}

func toRSAPublicKey(key jwk) (*rsa.PublicKey, error) {
	if key.KeyType != "RSA" {
		return nil, fmt.Errorf("unsupported jwk type %q", key.KeyType)
	}
	modulusBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		return nil, err
	}
	exponentBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		return nil, err
	}
	if len(exponentBytes) == 0 {
		return nil, errors.New("empty rsa exponent")
	}
	exponent := 0
	for _, b := range exponentBytes {
		exponent = exponent<<8 + int(b)
	}
	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(modulusBytes),
		E: exponent,
	}, nil
}

func HasAuthority(auth *Authentication, expected string) bool {
	if auth == nil {
		return false
	}
	for _, authority := range auth.Authorities {
		if authority == expected {
			return true
		}
	}
	return false
}
