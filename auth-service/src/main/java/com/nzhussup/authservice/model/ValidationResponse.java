package com.nzhussup.authservice.model;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.util.List;

@Data
@AllArgsConstructor
public class ValidationResponse {
    private String username;
    private List<String> roles;
}