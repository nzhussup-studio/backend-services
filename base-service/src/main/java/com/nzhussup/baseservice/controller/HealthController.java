package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping(AppConfig.baseApiPath + "health")
public class HealthController {

    @GetMapping
    public ResponseEntity<Map<String, String>> health() {
        return ResponseEntity.ok(Map.of("message", "OK"));
    }
}
