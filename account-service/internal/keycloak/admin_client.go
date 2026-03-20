package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	customerrors "account-service/internal/errors"
)

type AdminClientConfig struct {
	BaseURL           string
	Realm             string
	AdminRealm        string
	AdminClientID     string
	AdminClientSecret string
	HTTPTimeout       time.Duration
}

type AdminClient struct {
	config     AdminClientConfig
	httpClient *http.Client
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type userRepresentation struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewAdminClient(config AdminClientConfig) *AdminClient {
	return &AdminClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.HTTPTimeout,
		},
	}
}

func (c *AdminClient) LogoutUser(ctx context.Context, userID string) error {
	adminToken, err := c.adminAccessToken(ctx)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/admin/realms/%s/users/%s/logout", strings.TrimRight(c.config.BaseURL, "/"), url.PathEscape(c.config.Realm), url.PathEscape(userID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return customerrors.Wrap(customerrors.ErrInternal, "failed to create logout request")
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return customerrors.Wrap(customerrors.ErrInternal, "failed to call Keycloak logout endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return c.adminAPIError(resp, "failed to logout current user")
}

func (c *AdminClient) DeleteUser(ctx context.Context, userID string) error {
	adminToken, err := c.adminAccessToken(ctx)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/admin/realms/%s/users/%s", strings.TrimRight(c.config.BaseURL, "/"), url.PathEscape(c.config.Realm), url.PathEscape(userID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return customerrors.Wrap(customerrors.ErrInternal, "failed to create delete request")
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return customerrors.Wrap(customerrors.ErrInternal, "failed to call Keycloak delete endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return c.adminAPIError(resp, "failed to delete current user")
}

func (c *AdminClient) ResolveUserID(ctx context.Context, subject string, username string, email string) (string, error) {
	if subject != "" {
		return subject, nil
	}

	adminToken, err := c.adminAccessToken(ctx)
	if err != nil {
		return "", err
	}

	if username != "" {
		userID, err := c.findUserID(ctx, adminToken, "username", username)
		if err == nil {
			return userID, nil
		}
		if !strings.Contains(err.Error(), customerrors.ErrNotFound.Error()) {
			return "", err
		}
	}

	if email != "" {
		return c.findUserID(ctx, adminToken, "email", email)
	}

	return "", customerrors.Wrap(customerrors.ErrUnauthorized, "token does not contain a usable account identifier")
}

func (c *AdminClient) adminAccessToken(ctx context.Context) (string, error) {
	if c.config.AdminClientID == "" || c.config.AdminClientSecret == "" {
		return "", customerrors.Wrap(customerrors.ErrInternal, "Keycloak admin client credentials are not configured")
	}

	form := url.Values{}
	form.Set("client_id", c.config.AdminClientID)
	form.Set("grant_type", "client_credentials")
	form.Set("client_secret", c.config.AdminClientSecret)

	endpoint := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", strings.TrimRight(c.config.BaseURL, "/"), url.PathEscape(c.config.AdminRealm))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to create admin token request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to get Keycloak admin token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.adminAPIError(resp, "failed to get Keycloak admin token")
	}

	var payload tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to decode admin token response")
	}

	if payload.AccessToken == "" {
		return "", customerrors.Wrap(customerrors.ErrInternal, "Keycloak admin token response is empty")
	}

	return payload.AccessToken, nil
}

func (c *AdminClient) findUserID(ctx context.Context, adminToken string, field string, value string) (string, error) {
	query := url.Values{}
	query.Set("exact", "true")
	query.Set(field, value)

	endpoint := fmt.Sprintf(
		"%s/admin/realms/%s/users?%s",
		strings.TrimRight(c.config.BaseURL, "/"),
		url.PathEscape(c.config.Realm),
		query.Encode(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to create user lookup request")
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to call Keycloak user lookup endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", c.adminAPIError(resp, "failed to resolve current user")
	}

	var users []userRepresentation
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return "", customerrors.Wrap(customerrors.ErrInternal, "failed to decode Keycloak user lookup response")
	}

	for _, user := range users {
		if user.ID == "" {
			continue
		}
		if field == "username" && user.Username == value {
			return user.ID, nil
		}
		if field == "email" && strings.EqualFold(user.Email, value) {
			return user.ID, nil
		}
	}

	return "", customerrors.Wrap(customerrors.ErrNotFound, "current user not found")
}

func (c *AdminClient) adminAPIError(resp *http.Response, prefix string) error {
	body, _ := io.ReadAll(resp.Body)
	message := strings.TrimSpace(string(body))

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return customerrors.Wrap(customerrors.ErrUnauthorized, prefix)
	case http.StatusForbidden:
		return customerrors.Wrap(customerrors.ErrForbidden, prefix)
	case http.StatusNotFound:
		return customerrors.Wrap(customerrors.ErrNotFound, prefix)
	default:
		if message == "" {
			return customerrors.Wrap(customerrors.ErrInternal, prefix)
		}
		return customerrors.Wrap(customerrors.ErrInternal, fmt.Sprintf("%s: %s", prefix, message))
	}
}
