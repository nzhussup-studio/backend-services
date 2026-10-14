package com.nzhussup.baseservice.security;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpMethod;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.AnonymousAuthenticationFilter;

@Configuration
public class SecurityConfig {

    @Bean
    SecurityFilterChain securityFilterChain(HttpSecurity http, TrustedIdentityFilter trustedIdentityFilter)
            throws Exception {
        http
                .csrf(csrf -> csrf.disable())
                .authorizeHttpRequests(requests -> {
                    requests.requestMatchers("/error").permitAll();
                    requests.requestMatchers("/v1/base/cv-generator-preference/**").authenticated();
                    requests.requestMatchers(HttpMethod.GET).permitAll();
                    // Envoy authorizes admin routes before they reach this service.
                    requests.anyRequest().authenticated();
                })
                .sessionManagement(session ->
                        session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
                .addFilterBefore(trustedIdentityFilter, AnonymousAuthenticationFilter.class);

        return http.build();
    }
}
