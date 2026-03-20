package security

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	custom_errors "image-service/internal/errors"
	"io"
	"math/big"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ValidationResponse struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func Validate(c *gin.Context, config *AuthConfig, token string) (ValidationResponse, error) {
	claims, err := parseTokenClaims(token)
	if err != nil {
		return ValidationResponse{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Invalid token")
	}

	valid, err := verifyToken(c, config, token, claims)
	if err != nil {
		return ValidationResponse{}, custom_errors.NewError(custom_errors.ErrInternalServer, "Failed to validate token")
	}
	if !valid {
		return ValidationResponse{}, custom_errors.NewError(custom_errors.ErrUnauthorized, "Invalid token")
	}

	return ValidationResponse{
		Username: firstNonEmpty(claims.PreferredUsername, claims.Username, claims.Sub),
		Roles:    collectRoles(claims, config.BackendClientID),
	}, nil
}

func GetToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", custom_errors.NewError(custom_errors.ErrUnauthorized, "Authorization header missing or invalid")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func isAdmin(roles []string) bool {
	return slices.Contains(roles, "ROLE_ADMIN")
}

type tokenClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Username          string `json:"username"`
	Sub               string `json:"sub"`
	Iss               string `json:"iss"`
	Aud               any    `json:"aud"`
	Exp               int64  `json:"exp"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
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

func decodeBase64URL(value string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(value)
}

func verifyToken(c *gin.Context, config *AuthConfig, token string, claims *tokenClaims) (bool, error) {
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

	if claims.Exp <= 0 || claims.Exp <= unixNow() {
		return false, nil
	}

	keys, err := fetchJWKS(c, config.JWKSetURL)
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

func fetchJWKS(c *gin.Context, jwkSetURL string) (*jwks, error) {
	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, jwkSetURL, nil)
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

func unixNow() int64 {
	return timeNow().Unix()
}

var timeNow = func() time.Time {
	return time.Now()
}

func collectRoles(claims *tokenClaims, backendClientID string) []string {
	unique := map[string]struct{}{}
	roles := make([]string, 0)

	appendRole := func(role string) {
		if role == "" {
			return
		}
		if _, exists := unique[role]; exists {
			return
		}
		unique[role] = struct{}{}
		roles = append(roles, role)
	}

	for _, role := range claims.RealmAccess.Roles {
		appendRole(role)
	}

	if backendClientRoles, exists := claims.ResourceAccess[backendClientID]; exists {
		for _, role := range backendClientRoles.Roles {
			appendRole(role)
		}
	}

	return roles
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return "unknown"
}

func shouldAuthenticate(c *gin.Context, rules []AuthRule) bool {
	requestPath := c.Request.URL.Path

	for _, rule := range rules {
		if rule.Path != requestPath {
			continue
		}

		match := true
		for key, val := range rule.QueryParams {
			if c.Query(key) != val {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

var CheckIsAdmin = func(c *gin.Context, config *AuthConfig) error {
	token, err := GetToken(c)
	if err != nil {
		return err
	}
	validationResponse, err := Validate(c, config, token)
	if err != nil {
		return err
	}
	if isAdmin(validationResponse.Roles) {
		return nil
	}
	return custom_errors.NewError(custom_errors.ErrForbidden, "User is not an admin")
}
