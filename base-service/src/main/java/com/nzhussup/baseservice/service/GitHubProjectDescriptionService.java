package com.nzhussup.baseservice.service;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestClient;

import java.net.URI;
import java.util.Locale;
import java.util.Optional;

@Service
public class GitHubProjectDescriptionService {

    private static final Logger log = LoggerFactory.getLogger(GitHubProjectDescriptionService.class);
    private static final String GITHUB_API_VERSION = "2022-11-28";

    private final RestClient restClient;

    public GitHubProjectDescriptionService(
            RestClient.Builder restClientBuilder,
            @Value("${app.github.token:}") String githubToken) {
        RestClient.Builder builder = restClientBuilder
                .baseUrl("https://api.github.com")
                .defaultHeader(HttpHeaders.ACCEPT, "application/vnd.github+json")
                .defaultHeader("X-GitHub-Api-Version", GITHUB_API_VERSION);

        if (!githubToken.isBlank()) {
            builder.defaultHeader(HttpHeaders.AUTHORIZATION, "Bearer " + githubToken);
        }

        this.restClient = builder.build();
    }

    /**
     * Empty strings are deliberately returned (and cached) when no description can be
     * obtained. The project response converts them to null.
     */
    @Cacheable(value = "githubProjectDescriptions", key = "#repositoryUrl", sync = true)
    public String getDescription(String repositoryUrl) {
        Optional<RepositoryCoordinates> coordinates = RepositoryCoordinates.from(repositoryUrl);
        if (coordinates.isEmpty()) {
            return "";
        }

        RepositoryCoordinates repository = coordinates.get();
        try {
            GitHubRepository response = restClient.get()
                    .uri(uriBuilder -> uriBuilder
                            .pathSegment("repos", repository.owner(), repository.name())
                            .build())
                    .retrieve()
                    .body(GitHubRepository.class);
            return response == null || response.description() == null
                    ? ""
                    : response.description();
        } catch (Exception exception) {
            log.warn("Could not fetch GitHub description for {}/{}: {}",
                    repository.owner(), repository.name(), exception.getMessage());
            return "";
        }
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
