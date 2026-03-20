package com.nzhussup.baseservice.security;

import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.security.oauth2.core.OAuth2TokenValidator;
import org.springframework.security.oauth2.core.OAuth2TokenValidatorResult;
import org.springframework.security.oauth2.jwt.Jwt;

import java.util.List;

public class OptionalIssuerAudienceValidator implements OAuth2TokenValidator<Jwt> {

    private final String expectedIssuer;
    private final String expectedAudience;

    public OptionalIssuerAudienceValidator(String expectedIssuer, String expectedAudience) {
        this.expectedIssuer = expectedIssuer;
        this.expectedAudience = expectedAudience;
    }

    @Override
    public OAuth2TokenValidatorResult validate(Jwt token) {
        if (expectedIssuer != null && !expectedIssuer.isBlank()) {
            String issuer = token.getIssuer() != null ? token.getIssuer().toString() : null;
            if (!expectedIssuer.equals(issuer)) {
                return OAuth2TokenValidatorResult.failure(new OAuth2Error("invalid_token", "Invalid issuer", null));
            }
        }

        if (expectedAudience != null && !expectedAudience.isBlank()) {
            List<String> audiences = token.getAudience();
            if (audiences == null || !audiences.contains(expectedAudience)) {
                return OAuth2TokenValidatorResult.failure(new OAuth2Error("invalid_token", "Missing required audience", null));
            }
        }

        return OAuth2TokenValidatorResult.success();
    }
}
