package com.nzhussup.userservice.security;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.nzhussup.userservice.config.AppConfig;
import lombok.Getter;
import lombok.Setter;
import org.springframework.http.HttpHeaders;
import jakarta.servlet.Filter;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.stream.Collectors;

@Component
public class AuthServiceFilter implements Filter {

    private static final ObjectMapper objectMapper = new ObjectMapper();

    private static final String AUTH_SERVICE_URL = System.getenv("AUTH_SERVICE_URL");

    private final HttpClient httpClient;

    public AuthServiceFilter() {
        if (AUTH_SERVICE_URL == null || AUTH_SERVICE_URL.isEmpty()) {
            throw new IllegalStateException("AUTH_SERVICE_URL environment variable is not set");
        }
        this.httpClient = HttpClient.newHttpClient();
    }

    @Override
    public void doFilter(
            jakarta.servlet.ServletRequest servletRequest,
            jakarta.servlet.ServletResponse servletResponse,
            FilterChain filterChain) throws IOException, ServletException {

        HttpServletRequest request = (HttpServletRequest) servletRequest;
        HttpServletResponse response = (HttpServletResponse) servletResponse;

        String authHeader = request.getHeader(HttpHeaders.AUTHORIZATION);

        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            filterChain.doFilter(request, response);
            return;
        }

        String token = authHeader.substring(7);

        try {
            ValidationResponse validationResponse = validateToken(token);

            SecurityContextHolder.getContext().setAuthentication(
                    new UsernamePasswordAuthenticationToken(
                            validationResponse.getUsername(),
                            null,
                            validationResponse.getAuthorities()));

            filterChain.doFilter(request, response);

        } catch (Exception e) {
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            response.getWriter().write("Unauthorized: " + e.getMessage());
        }
    }

    private ValidationResponse validateToken(String token) throws IOException, InterruptedException {
        ValidationRequest validationRequest = new ValidationRequest(token);
        String requestBody = objectMapper.writeValueAsString(validationRequest);

        HttpRequest httpRequest = HttpRequest.newBuilder()
                .uri(URI.create(AUTH_SERVICE_URL + AppConfig.authValidationPath))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(requestBody, StandardCharsets.UTF_8))
                .build();

        HttpResponse<String> httpResponse = httpClient.send(httpRequest, HttpResponse.BodyHandlers.ofString());

        if (httpResponse.statusCode() != 200) {
            throw new IOException("Auth service returned HTTP " + httpResponse.statusCode());
        }

        return objectMapper.readValue(httpResponse.body(), ValidationResponse.class);
    }

    @Getter
    @Setter
    public static class ValidationRequest {
        private String token;

        public ValidationRequest(String token) {
            this.token = token;
        }
    }

    @Getter
    @Setter
    public static class ValidationResponse {
        private String username;
        private List<String> roles;

        public List<SimpleGrantedAuthority> getAuthorities() {
            return roles.stream()
                    .map(SimpleGrantedAuthority::new)
                    .collect(Collectors.toList());
        }
    }
}