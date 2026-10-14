package auth

type AuthRule struct {
	Path        string
	QueryParams map[string]string
}

type AuthConfig struct {
	Rules []AuthRule
}
