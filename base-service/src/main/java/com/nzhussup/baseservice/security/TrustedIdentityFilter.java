package com.nzhussup.baseservice.security;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

@Component
public class TrustedIdentityFilter extends OncePerRequestFilter {

    private static final String USER_ID_HEADER = "X-User-ID";
    private static final String USER_NAME_HEADER = "X-User-Name";
    private static final String USER_ROLES_HEADER = "X-User-Roles";

    private final ObjectMapper objectMapper;

    public TrustedIdentityFilter(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
    }

    @Override
    protected void doFilterInternal(
            HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain) throws ServletException, IOException {
        String subject = trimToNull(request.getHeader(USER_ID_HEADER));
        if (subject != null) {
            String username = trimToNull(request.getHeader(USER_NAME_HEADER));
            List<SimpleGrantedAuthority> authorities = decodeRoles(request.getHeader(USER_ROLES_HEADER)).stream()
                    .map(SimpleGrantedAuthority::new)
                    .toList();
            var authentication = UsernamePasswordAuthenticationToken.authenticated(
                    username != null ? username : subject,
                    null,
                    authorities);
            SecurityContextHolder.getContext().setAuthentication(authentication);
        }

        filterChain.doFilter(request, response);
    }

    private Set<String> decodeRoles(String value) {
        String rolesHeader = trimToNull(value);
        if (rolesHeader == null) {
            return Set.of();
        }

        for (Base64.Decoder decoder : List.of(
                Base64.getDecoder(), Base64.getUrlDecoder(),
                Base64.getMimeDecoder())) {
            try {
                byte[] decoded = decoder.decode(rolesHeader);
                List<String> roles = objectMapper.readValue(decoded, new TypeReference<ArrayList<String>>() {});
                return normalize(roles);
            } catch (IllegalArgumentException | IOException ignored) {
                // Try the next supported Envoy base64 representation.
            }
        }

        if (!rolesHeader.matches(".*[\\[\\]{}\"].*")) {
            return normalize(List.of(rolesHeader.split(",")));
        }

        return Set.of();
    }

    private Set<String> normalize(List<String> roles) {
        return roles.stream()
                .map(String::trim)
                .filter(role -> !role.isEmpty())
                .collect(Collectors.toUnmodifiableSet());
    }

    private String trimToNull(String value) {
        if (value == null || value.isBlank()) {
            return null;
        }
        return value.trim();
    }
}
