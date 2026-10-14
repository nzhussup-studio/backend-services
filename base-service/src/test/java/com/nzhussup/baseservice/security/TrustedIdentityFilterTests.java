package com.nzhussup.baseservice.security;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockFilterChain;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;
import org.springframework.security.core.context.SecurityContextHolder;

class TrustedIdentityFilterTests {

    private final TrustedIdentityFilter filter = new TrustedIdentityFilter();

    @AfterEach
    void clearSecurityContext() {
        SecurityContextHolder.clearContext();
    }

    @Test
    void createsAuthenticationFromTrustedHeaders() throws Exception {
        var request = new MockHttpServletRequest();
        request.addHeader("X-User-ID", "user-123");
        request.addHeader("X-User-Name", "admin");

        filter.doFilter(request, new MockHttpServletResponse(), new MockFilterChain());

        var authentication = SecurityContextHolder.getContext().getAuthentication();
        assertThat(authentication).isNotNull();
        assertThat(authentication.getName()).isEqualTo("admin");
        assertThat(authentication.isAuthenticated()).isTrue();
        assertThat(authentication.getAuthorities()).isEmpty();
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
