package main

import "image-service/internal/config/security"

func GetSecurityConfig(config *config) *security.AuthConfig {
	return &security.AuthConfig{
		JWKSetURL:        config.keycloakConfig.jwkSetURL,
		ExpectedIssuer:   config.keycloakConfig.expectedIssuer,
		ExpectedAudience: config.keycloakConfig.expectedAudience,
		BackendClientID:  config.keycloakConfig.backendClientID,
		Rules: []security.AuthRule{
			{
				Path: config.apiBasePath,
				QueryParams: map[string]string{
					"type": "private",
				},
			},
			{
				Path: config.apiBasePath,
				QueryParams: map[string]string{
					"type": "semi-public",
				},
			},
			{
				Path: config.apiBasePath,
				QueryParams: map[string]string{
					"type": "all",
				},
			},
		},
	}
}
