package auth

// ProtectedReadRule identifies a GET endpoint that needs a trusted identity.
type ProtectedReadRule struct {
	Path        string
	QueryParams map[string]string
}

type ProtectedReadRules struct {
	Rules []ProtectedReadRule
}
