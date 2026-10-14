package com.nzhussup.baseservice.service;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.http.HttpHeaders;
import org.springframework.http.client.JdkClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;
import org.springframework.web.client.RestClientResponseException;

import java.net.URI;
import java.net.http.HttpClient;
import java.util.Locale;
import java.util.Optional;

@Service
public class GitHubProjectDescriptionService {

    private static final Logger log = LoggerFactory.getLogger(GitHubProjectDescriptionService.class);
    private static final String GITHUB_API_VERSION = "2022-11-28";
    private static final String USER_AGENT = "nzhussup-studio-base-service";

    private final RestClient restClient;
    private final RestClient anonymousRestClient;
    private final boolean authenticated;

    @Autowired
    public GitHubProjectDescriptionService(
            RestClient.Builder restClientBuilder,
            @Value("${app.github.token:}") String githubToken) {
        this.restClient = createClient(restClientBuilder.clone(), githubToken);
        this.anonymousRestClient = githubToken.isBlank()
                ? this.restClient
                : createClient(restClientBuilder.clone(), "");
        this.authenticated = !githubToken.isBlank();
    }

    private static RestClient createClient(RestClient.Builder builder, String githubToken) {
        builder
                .baseUrl("https://api.github.com")
                .requestFactory(new JdkClientHttpRequestFactory(HttpClient.newBuilder()
                        .followRedirects(HttpClient.Redirect.NORMAL)
                        .build()))
                .defaultHeader(HttpHeaders.ACCEPT, "application/vnd.github+json")
                .defaultHeader(HttpHeaders.USER_AGENT, USER_AGENT)
                .defaultHeader("X-GitHub-Api-Version", GITHUB_API_VERSION);

        if (!githubToken.isBlank()) {
            builder.defaultHeader(HttpHeaders.AUTHORIZATION, "Bearer " + githubToken);
        }
        return builder.build();
    }

    GitHubProjectDescriptionService(RestClient restClient) {
        this.restClient = restClient;
        this.anonymousRestClient = restClient;
        this.authenticated = false;
    }

    GitHubProjectDescriptionService(RestClient restClient, RestClient anonymousRestClient) {
        this.restClient = restClient;
        this.anonymousRestClient = anonymousRestClient;
        this.authenticated = true;
    }

    /**
     * Empty strings are not cached so a transient GitHub failure can recover on the
     * next request. The project response converts them to null.
     */
    @Cacheable(value = "githubProjectDescriptions", key = "#repositoryUrl",
            unless = "#result == null || #result.isBlank()")
    public String getDescription(String repositoryUrl) {
        Optional<RepositoryCoordinates> coordinates = RepositoryCoordinates.from(repositoryUrl);
        if (coordinates.isEmpty()) {
            return "";
        }

        RepositoryCoordinates repository = coordinates.get();
        try {
            return fetchDescription(restClient, repository);
        } catch (RestClientResponseException exception) {
            if (authenticated && (exception.getStatusCode().value() == 401
                    || exception.getStatusCode().value() == 403)) {
                log.warn("Authenticated GitHub request for {}/{} returned {}; retrying publicly",
                        repository.owner(), repository.name(), exception.getStatusCode().value());
                try {
                    return fetchDescription(anonymousRestClient, repository);
                } catch (Exception fallbackException) {
                    log.warn("Could not fetch public GitHub description for {}/{}: {}",
                            repository.owner(), repository.name(), fallbackException.getMessage());
                    return "";
                }
            }
            log.warn("Could not fetch GitHub description for {}/{} (HTTP {}): {}",
                    repository.owner(), repository.name(), exception.getStatusCode().value(),
                    exception.getMessage());
            return "";
        } catch (Exception exception) {
            log.warn("Could not fetch GitHub description for {}/{}: {}",
                    repository.owner(), repository.name(), exception.getMessage());
            return "";
        }
    }

    private String fetchDescription(RestClient client, RepositoryCoordinates repository) {
        GitHubRepository response = client.get()
                    .uri(uriBuilder -> uriBuilder
                            .pathSegment("repos", repository.owner(), repository.name())
                            .build())
                    .retrieve()
                    .body(GitHubRepository.class);
        return response == null || response.description() == null
                ? ""
                : response.description();
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    private record GitHubRepository(String description) {
    }

    private record RepositoryCoordinates(String owner, String name) {

        private static Optional<RepositoryCoordinates> from(String repositoryUrl) {
            if (repositoryUrl == null || repositoryUrl.isBlank()) {
                return Optional.empty();
            }

            try {
                URI uri = URI.create(repositoryUrl.trim());
                if (uri.getHost() == null
                        || !uri.getHost().toLowerCase(Locale.ROOT).equals("github.com")) {
                    return Optional.empty();
                }

                String[] parts = uri.getPath().split("/");
                if (parts.length != 3 || parts[1].isBlank() || parts[2].isBlank()) {
                    return Optional.empty();
                }

                String repositoryName = parts[2].endsWith(".git")
                        ? parts[2].substring(0, parts[2].length() - 4)
                        : parts[2];
                return repositoryName.isBlank()
                        ? Optional.empty()
                        : Optional.of(new RepositoryCoordinates(parts[1], repositoryName));
            } catch (IllegalArgumentException exception) {
                return Optional.empty();
            }
        }
    }
}
