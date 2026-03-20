package com.nzhussup.baseservice.security;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.oauth2.core.DelegatingOAuth2TokenValidator;
import org.springframework.security.oauth2.core.OAuth2TokenValidator;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.JwtValidators;
import org.springframework.security.oauth2.jwt.NimbusJwtDecoder;
import org.springframework.security.web.SecurityFilterChain;

@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    SecurityFilterChain securityFilterChain(
            HttpSecurity http,
            KeycloakJwtAuthenticationConverter keycloakJwtAuthenticationConverter)
            throws Exception {

        http
                .csrf(csfr -> csfr.disable())
                .authorizeHttpRequests(requests -> {
                    requests.requestMatchers(HttpMethod.GET).permitAll();
                    requests.anyRequest().hasRole("ADMIN");
                })
                .sessionManagement(session -> {
                    session.sessionCreationPolicy(SessionCreationPolicy.STATELESS);
                })
                .oauth2ResourceServer(oauth2 -> oauth2.jwt(
                        jwt -> jwt.jwtAuthenticationConverter(keycloakJwtAuthenticationConverter)));

        return http.build();
    }

    @Bean
    JwtDecoder jwtDecoder(
            @Value("${app.security.keycloak.jwk-set-uri}") String jwkSetUri,
            @Value("${app.security.keycloak.expected-issuer:}") String expectedIssuer,
            @Value("${app.security.keycloak.expected-audience:}") String expectedAudience) {
        NimbusJwtDecoder jwtDecoder = NimbusJwtDecoder.withJwkSetUri(jwkSetUri).build();

        OAuth2TokenValidator<Jwt> defaultValidators = JwtValidators.createDefault();
        OAuth2TokenValidator<Jwt> issuerAudienceValidator = new OptionalIssuerAudienceValidator(
                expectedIssuer,
                expectedAudience
        );
        jwtDecoder.setJwtValidator(new DelegatingOAuth2TokenValidator<>(defaultValidators, issuerAudienceValidator));

        return jwtDecoder;
    }

    @Bean
    KeycloakJwtAuthenticationConverter keycloakJwtAuthenticationConverter(
            @Value("${app.security.keycloak.backend-client-id}") String backendClientId) {
        return new KeycloakJwtAuthenticationConverter(backendClientId);
    }
}
