package security

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	customerrors "account-service/internal/errors"

	"github.com/gin-gonic/gin"
)

type ValidationResponse struct {
	Subject  string `json:"subject"`
	Username string `json:"username"`
}

type tokenClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Username          string `json:"username"`
	Sub               string `json:"sub"`
	Iss               string `json:"iss"`
	Aud               any    `json:"aud"`
	Exp               int64  `json:"exp"`
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func Validate(c *gin.Context, config *AuthConfig, token string) (ValidationResponse, error) {
	claims, err := parseTokenClaims(token)
	if err != nil {
		return ValidationResponse{}, customerrors.Wrap(customerrors.ErrUnauthorized, "invalid token")
	}

	valid, err := verifyToken(c.Request.Context(), config, token, claims)
	if err != nil {
		return ValidationResponse{}, customerrors.Wrap(customerrors.ErrInternal, "failed to validate token")
	}

	if !valid {
		return ValidationResponse{}, customerrors.Wrap(customerrors.ErrUnauthorized, "invalid token")
	}

	return ValidationResponse{
		Subject:  claims.Sub,
		Username: firstNonEmpty(claims.PreferredUsername, claims.Username, claims.Sub),
	}, nil
}

func GetToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", customerrors.Wrap(customerrors.ErrUnauthorized, "authorization header missing or invalid")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", customerrors.Wrap(customerrors.ErrUnauthorized, "authorization header missing or invalid")
	}

	return token, nil
}

func parseTokenClaims(token string) (*tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	payload, err := decodeBase64URL(parts[1])
	if err != nil {
		return nil, err
	}

	var claims tokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func verifyToken(ctx context.Context, config *AuthConfig, token string, claims *tokenClaims) (bool, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false, nil
	}

	headerBytes, err := decodeBase64URL(parts[0])
	if err != nil {
		return false, nil
	}

	var header tokenHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return false, nil
	}

	if header.Alg != "RS256" || header.Kid == "" {
		return false, nil
	}

	if config.ExpectedIssuer != "" && claims.Iss != config.ExpectedIssuer {
		return false, nil
	}

	if config.ExpectedAudience != "" && !hasAudience(claims.Aud, config.ExpectedAudience) {
		return false, nil
	}

	if claims.Sub == "" || claims.Exp <= 0 || claims.Exp <= time.Now().Unix() {
		return false, nil
	}

	keys, err := fetchJWKS(ctx, config.JWKSetURL)
	if err != nil {
		return false, err
	}

	publicKey, err := findPublicKey(keys, header.Kid)
	if err != nil {
		return false, nil
	}

	signingInput := parts[0] + "." + parts[1]
	signature, err := decodeBase64URL(parts[2])
	if err != nil {
		return false, nil
	}

	hashed := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature); err != nil {
		return false, nil
	}

	return true, nil
}

func fetchJWKS(ctx context.Context, jwkSetURL string) (*jwks, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, jwkSetURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed jwks
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	return &parsed, nil
}

func findPublicKey(keys *jwks, keyID string) (*rsa.PublicKey, error) {
	for _, key := range keys.Keys {
		if key.Kid != keyID || key.Kty != "RSA" || key.Use != "sig" {
			continue
		}

		modulusBytes, err := decodeBase64URL(key.N)
		if err != nil {
			return nil, err
		}

		exponentBytes, err := decodeBase64URL(key.E)
		if err != nil {
			return nil, err
		}

		modulus := new(big.Int).SetBytes(modulusBytes)
		exponent := new(big.Int).SetBytes(exponentBytes)

		return &rsa.PublicKey{
			N: modulus,
			E: int(exponent.Int64()),
		}, nil
	}

	return nil, fmt.Errorf("signing key not found")
}

func decodeBase64URL(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}

func hasAudience(audClaim any, expectedAudience string) bool {
	switch aud := audClaim.(type) {
	case string:
		return aud == expectedAudience
	case []any:
		for _, entry := range aud {
			if audString, ok := entry.(string); ok && audString == expectedAudience {
				return true
			}
		}
	}

	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return "unknown"
}
