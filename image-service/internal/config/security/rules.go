package security

type AuthRule struct {
	Path        string
	QueryParams map[string]string
}

type AuthConfig struct {
	AuthServiceURL string
	ValidationURL  string
	Rules          []AuthRule
}
