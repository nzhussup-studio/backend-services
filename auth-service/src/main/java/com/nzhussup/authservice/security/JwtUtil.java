package com.nzhussup.authservice.security;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.exceptions.JWTVerificationException;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.interfaces.JWTVerifier;
import com.nzhussup.authservice.config.AppConfig;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.stereotype.Component;

import java.util.Collection;
import java.util.Date;
import java.util.stream.Collectors;

@Component
public class JwtUtil {

    @Value("security.jwt.secret-key")
    private String SECRET_KEY;

    public String generateToken(String username, Collection<? extends GrantedAuthority> roles) {

        String rolesString = roles.stream()
                .map(GrantedAuthority::getAuthority)
                .collect(Collectors.joining(","));

        return JWT.create()
                .withSubject(username)
                .withClaim("roles", rolesString)
                .withIssuedAt(new Date())
                .withExpiresAt(new Date(System.currentTimeMillis() + AppConfig.EXPIRATION_TIME))
                .sign(Algorithm.HMAC256(SECRET_KEY));
    }

    public long getExpirationTime(String token) {
        return JWT.decode(token).getClaim("exp").asLong();
    }

    public DecodedJWT validateToken(String token) throws JWTVerificationException {
        JWTVerifier verifier = JWT.require(Algorithm.HMAC256(SECRET_KEY))
                .build();
        return verifier.verify(token);
    }

    public String getUsernameFromToken(String token) {
        return validateToken(token).getSubject();
    }
}