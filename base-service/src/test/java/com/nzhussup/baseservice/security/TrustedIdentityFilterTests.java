package com.nzhussup.baseservice.security;

import static org.assertj.core.api.Assertions.assertThat;

import com.fasterxml.jackson.databind.ObjectMapper;
import java.nio.charset.StandardCharsets;
import java.util.Base64;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockFilterChain;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;
import org.springframework.security.core.context.SecurityContextHolder;

class TrustedIdentityFilterTests {

    private final TrustedIdentityFilter filter = new TrustedIdentityFilter(new ObjectMapper());

    @AfterEach
    void clearSecurityContext() {
        SecurityContextHolder.clearContext();
    }

    @Test
    void createsAuthenticationFromTrustedHeaders() throws Exception {
        var request = new MockHttpServletRequest();
        request.addHeader("X-User-ID", "user-123");
        request.addHeader("X-User-Name", "admin");
        request.addHeader("X-User-Roles", Base64.getEncoder().encodeToString(
                "[\"ROLE_USER\",\"ROLE_ADMIN\"]".getBytes(StandardCharsets.UTF_8)));

        filter.doFilter(request, new MockHttpServletResponse(), new MockFilterChain());

        var authentication = SecurityContextHolder.getContext().getAuthentication();
        assertThat(authentication).isNotNull();
        assertThat(authentication.getName()).isEqualTo("admin");
        assertThat(authentication.getAuthorities())
                .extracting(Object::toString)
                .contains("ROLE_ADMIN", "ROLE_USER");
    }

    @Test
    void remainsAnonymousWhenTrustedSubjectIsMissing() throws Exception {
        filter.doFilter(
                new MockHttpServletRequest(),
                new MockHttpServletResponse(),
                new MockFilterChain());

        assertThat(SecurityContextHolder.getContext().getAuthentication()).isNull();
    }
}
