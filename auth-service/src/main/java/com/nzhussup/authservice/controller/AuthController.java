package com.nzhussup.authservice.controller;

import com.auth0.jwt.exceptions.JWTVerificationException;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.nzhussup.authservice.model.AuthRequest;
import com.nzhussup.authservice.model.AuthResponse;
import com.nzhussup.authservice.model.ValidationRequest;
import com.nzhussup.authservice.model.ValidationResponse;
import com.nzhussup.authservice.security.JwtUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequiredArgsConstructor
@RequestMapping("/auth")
public class AuthController {

    private final BCryptPasswordEncoder bCryptPasswordEncoder = new BCryptPasswordEncoder(12);
    private final AuthenticationManager authenticationManager;
    private final JwtUtil jwtUtil;
    private final UserDetailsService userDetailsService;

    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody AuthRequest authRequest) {
        try {
            UserDetails userDetails = userDetailsService.loadUserByUsername(authRequest.getUsername());
            if (!bCryptPasswordEncoder.matches(authRequest.getPassword(), userDetails.getPassword())) {
                throw new BadCredentialsException("Bad credentials");
            }
            Authentication authentication = authenticationManager.authenticate(
                    new UsernamePasswordAuthenticationToken(authRequest.getUsername(), authRequest.getPassword()));
            String token = jwtUtil.generateToken(authRequest.getUsername(), userDetails.getAuthorities());
            Long expiration = jwtUtil.getExpirationTime(token);
            return ResponseEntity.ok(new AuthResponse(token, expiration));
        } catch (BadCredentialsException | UsernameNotFoundException e) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(e.getMessage());
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(e.getMessage());
        }
    }

    @PostMapping("/validate")
    public ResponseEntity<?> validateToken(@RequestBody ValidationRequest validationRequest) {
        try {
            DecodedJWT decodedJWT = jwtUtil.validateToken(validationRequest.getToken());

            String username = decodedJWT.getSubject();
            Object rolesClaim = decodedJWT.getClaim("roles").as(Object.class);
            List<String> roles;

            if (rolesClaim instanceof String) {
                roles = List.of(((String) rolesClaim).split(",")).stream()
                        .map(String::trim) // Remove any extra whitespace
                        .toList();
            } else if (rolesClaim instanceof List) {
                roles = decodedJWT.getClaim("roles").asList(String.class);
            } else {
                roles = List.of();
            }

            return ResponseEntity.ok(new ValidationResponse(username, roles));

        } catch (JWTVerificationException e) {
            return ResponseEntity.status(401).body("Invalid or expired token: " + e.getMessage());
        }
    }

}