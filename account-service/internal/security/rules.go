package security

type AuthConfig struct {
	JWKSetURL        string
	ExpectedIssuer   string
	ExpectedAudience string
}
