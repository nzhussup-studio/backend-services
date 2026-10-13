package com.nzhussup.baseservice.controller;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class HealthControllerTests {

    @Test
    void returnsOk() {
        var response = new HealthController().health();

        assertThat(response.getStatusCode().is2xxSuccessful()).isTrue();
        assertThat(response.getBody()).containsEntry("message", "OK");
    }
}
