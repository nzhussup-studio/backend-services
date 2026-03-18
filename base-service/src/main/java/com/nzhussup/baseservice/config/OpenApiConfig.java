package com.nzhussup.baseservice.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.Operation;
import io.swagger.v3.oas.models.PathItem;
import org.springdoc.core.customizers.OpenApiCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.LinkedHashMap;
import java.util.Locale;
import java.util.Map;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenApiCustomizer baseServiceOperationIdCustomizer() {
        return openApi -> {
            if (openApi.getPaths() == null) {
                return;
            }

            openApi.getPaths().forEach((path, pathItem) -> customizePathOperations(path, pathItem));
        };
    }

    private void customizePathOperations(String path, PathItem pathItem) {
        String resourceName = extractResourceName(path);
        if (resourceName == null) {
            return;
        }

        Map<PathItem.HttpMethod, String> operationNames = new LinkedHashMap<>();
        operationNames.put(PathItem.HttpMethod.GET, "list");
        operationNames.put(PathItem.HttpMethod.POST, "create");
        operationNames.put(PathItem.HttpMethod.PUT, "update");
        operationNames.put(PathItem.HttpMethod.DELETE, "delete");

        operationNames.forEach((method, prefix) -> {
            Operation operation = pathItem.readOperationsMap().get(method);
            if (operation == null) {
                return;
            }

            operation.setOperationId(prefix + resourceName);
        });
    }

    private String extractResourceName(String path) {
        String basePath = "/" + AppConfig.baseApiPath;
        if (!path.startsWith(basePath)) {
            return null;
        }

        String resourcePath = path.substring(basePath.length());
        int slashIndex = resourcePath.indexOf('/');
        if (slashIndex >= 0) {
            resourcePath = resourcePath.substring(0, slashIndex);
        }

        if (resourcePath.isBlank()) {
            return null;
        }

        return toPascalCase(resourcePath);
    }

    private String toPascalCase(String value) {
        StringBuilder builder = new StringBuilder();

        for (String part : value.split("-")) {
            if (part.isBlank()) {
                continue;
            }

            String normalized = part.toLowerCase(Locale.ROOT);
            builder.append(Character.toUpperCase(normalized.charAt(0)));
            builder.append(normalized.substring(1));
        }

        return builder.toString();
    }
}
