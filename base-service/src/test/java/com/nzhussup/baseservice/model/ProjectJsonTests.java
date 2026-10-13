package com.nzhussup.baseservice.model;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class ProjectJsonTests {

    private final ObjectMapper objectMapper = new ObjectMapper();

    @Test
    void descriptionIsReturnedButCannotBeSetByApiInput() throws Exception {
        Project project = Project.builder().description("From GitHub").build();

        assertThat(objectMapper.writeValueAsString(project))
                .contains("\"description\":\"From GitHub\"");

        Project submittedProject = objectMapper.readValue(
                "{\"name\":\"Project\",\"url\":\"https://github.com/example/project\","
                        + "\"description\":\"User supplied\"}",
                Project.class);
        assertThat(submittedProject.getDescription()).isNull();
    }
}
