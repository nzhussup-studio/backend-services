package com.nzhussup.baseservice.security;

import org.springframework.core.convert.converter.Converter;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;

import java.util.Collection;
import java.util.LinkedHashSet;
import java.util.Map;
import java.util.Set;

public class KeycloakJwtAuthenticationConverter implements Converter<Jwt, AbstractAuthenticationToken> {

    private final String backendClientId;

    public KeycloakJwtAuthenticationConverter(String backendClientId) {
        this.backendClientId = backendClientId;
    }

    @Override
    public AbstractAuthenticationToken convert(Jwt jwt) {
        Set<GrantedAuthority> authorities = new LinkedHashSet<>();
        addRoles(authorities, extractRealmRoles(jwt));
        addRoles(authorities, extractClientRoles(jwt, backendClientId));

        String principalName = jwt.getClaimAsString("preferred_username");
        if (principalName == null || principalName.isBlank()) {
            principalName = jwt.getSubject();
        }

        return new JwtAuthenticationToken(jwt, authorities, principalName);
    }

    @SuppressWarnings("unchecked")
    private Collection<String> extractRealmRoles(Jwt jwt) {
        Object realmAccess = jwt.getClaim("realm_access");
        if (!(realmAccess instanceof Map<?, ?> realmAccessMap)) {
            return Set.of();
        }

        Object roles = realmAccessMap.get("roles");
        if (!(roles instanceof Collection<?> roleCollection)) {
            return Set.of();
        }

        Set<String> resolvedRoles = new LinkedHashSet<>();
        for (Object role : roleCollection) {
            if (role instanceof String roleName && !roleName.isBlank()) {
                resolvedRoles.add(roleName);
            }
        }
        return resolvedRoles;
    }

    @SuppressWarnings("unchecked")
    private Collection<String> extractClientRoles(Jwt jwt, String clientId) {
        Object resourceAccess = jwt.getClaim("resource_access");
        if (!(resourceAccess instanceof Map<?, ?> resourceAccessMap)) {
            return Set.of();
        }

        Object clientEntry = resourceAccessMap.get(clientId);
        if (!(clientEntry instanceof Map<?, ?> clientMap)) {
            return Set.of();
        }

        Object roles = clientMap.get("roles");
        if (!(roles instanceof Collection<?> roleCollection)) {
            return Set.of();
        }

        Set<String> resolvedRoles = new LinkedHashSet<>();
        for (Object role : roleCollection) {
            if (role instanceof String roleName && !roleName.isBlank()) {
                resolvedRoles.add(roleName);
            }
        }
        return resolvedRoles;
    }

    private void addRoles(Set<GrantedAuthority> authorities, Collection<String> roles) {
        for (String role : roles) {
            authorities.add(new SimpleGrantedAuthority(role));
        }
    }
}
