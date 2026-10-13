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

class GitHubProjectDescriptionServiceTests {

    @Test
    void fetchesDescriptionFromGitHubRepositoryUrl() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder, "");

        server.expect(once(), requestTo("https://api.github.com/repos/example/project"))
                .andExpect(method(HttpMethod.GET))
                .andRespond(withSuccess("{\"description\":\"From GitHub\"}", MediaType.APPLICATION_JSON));

        assertThat(service.getDescription("https://github.com/example/project"))
                .isEqualTo("From GitHub");
        server.verify();
    }

    @Test
    void returnsEmptyDescriptionWhenGitHubHasNoDescription() {
        RestClient.Builder builder = RestClient.builder();
        MockRestServiceServer server = MockRestServiceServer.bindTo(builder).build();
        GitHubProjectDescriptionService service =
                new GitHubProjectDescriptionService(builder, "");

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
                new GitHubProjectDescriptionService(builder, "");

        assertThat(service.getDescription("https://example.com/example/project")).isEmpty();
        server.verify();
    }
}
