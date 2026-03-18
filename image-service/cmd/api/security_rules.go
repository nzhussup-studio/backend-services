package main

import "image-service/internal/config/security"

func GetSecurityConfig(config *config) *security.AuthConfig {
	return &security.AuthConfig{
		AuthServiceURL: config.apiGatewayURL,
		ValidationURL:  "/auth/validate",
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
