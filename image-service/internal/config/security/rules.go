package security

type AuthRule struct {
	Path        string
	QueryParams map[string]string
}

type AuthConfig struct {
	JWKSetURL        string
	ExpectedIssuer   string
	ExpectedAudience string
	BackendClientID  string
	Rules            []AuthRule
}
