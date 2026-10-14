package com.nzhussup.baseservice.service;

import org.junit.jupiter.api.Test;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.test.web.client.MockRestServiceServer;
import org.springframework.web.client.RestClient;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.client.ExpectedCount.once;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.method;
import static org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess;
import static org.springframework.test.web.client.response.MockRestResponseCreators.withUnauthorizedRequest;

class GitHubProjectDescriptionServiceTests {

    @Test
    void fetchesDescriptionFromGitHubRepositoryUrl() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder.baseUrl("https://api.github.com").build());

        server.expect(once(), requestTo("https://api.github.com/repos/example/project"))
                .andExpect(method(HttpMethod.GET))
                .andRespond(withSuccess("{\"description\":\"From GitHub\"}", MediaType.APPLICATION_JSON));

        assertThat(service.getDescription("https://github.com/example/project"))
                .isEqualTo("From GitHub");
        server.verify();
    }

    @Test
    void fetchesDescriptionForOrganizationWithHyphen() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder.baseUrl("https://api.github.com").build());

        server.expect(once(), requestTo(
                        "https://api.github.com/repos/nzhussup-studio/web-admin-panel"))
                .andExpect(method(HttpMethod.GET))
                .andRespond(withSuccess("{\"description\":\"Admin panel\"}", MediaType.APPLICATION_JSON));

        assertThat(service.getDescription(
                "https://github.com/nzhussup-studio/web-admin-panel"))
                .isEqualTo("Admin panel");
        server.verify();
    }

    @Test
    void retriesPubliclyWhenConfiguredTokenIsRejected() {
        RestClient.Builder authenticatedBuilder = RestClient.builder();
        MockRestServiceServer authenticatedServer =
                MockRestServiceServer.bindTo(authenticatedBuilder).build();
        RestClient.Builder anonymousBuilder = RestClient.builder();
        MockRestServiceServer anonymousServer = MockRestServiceServer.bindTo(anonymousBuilder).build();
        GitHubProjectDescriptionService service = new GitHubProjectDescriptionService(
                authenticatedBuilder.baseUrl("https://api.github.com").build(),
                anonymousBuilder.baseUrl("https://api.github.com").build());

        authenticatedServer.expect(once(), requestTo(
                        "https://api.github.com/repos/nzhussup-studio/web-admin-panel"))
                .andRespond(withUnauthorizedRequest());
        anonymousServer.expect(once(), requestTo(
                        "https://api.github.com/repos/nzhussup-studio/web-admin-panel"))
                .andRespond(withSuccess("{\"description\":\"Admin panel\"}", MediaType.APPLICATION_JSON));

        assertThat(service.getDescription(
                "https://github.com/nzhussup-studio/web-admin-panel"))
                .isEqualTo("Admin panel");
        authenticatedServer.verify();
        anonymousServer.verify();
    }

    @Test
    void returnsEmptyDescriptionWhenGitHubHasNoDescription() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder.baseUrl("https://api.github.com").build());

        server.expect(once(), requestTo("https://api.github.com/repos/example/project"))
                .andRespond(withSuccess("{\"description\":null}", MediaType.APPLICATION_JSON));

        assertThat(service.getDescription("https://github.com/example/project.git")).isEmpty();
        server.verify();
    }

    @Test
    void doesNotCallGitHubForNonGitHubUrl() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder.baseUrl("https://api.github.com").build());

        assertThat(service.getDescription("https://example.com/example/project")).isEmpty();
        server.verify();
    }
}
