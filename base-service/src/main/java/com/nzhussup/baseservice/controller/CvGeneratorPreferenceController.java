package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.CvGeneratorPreference;
import com.nzhussup.baseservice.service.CvGeneratorPreferenceService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping(AppConfig.adminApiPath + "cv-generator-preferences")
public class CvGeneratorPreferenceController {

    private final CvGeneratorPreferenceService preferenceService;

    @Autowired
    public CvGeneratorPreferenceController(CvGeneratorPreferenceService preferenceService) {
        this.preferenceService = preferenceService;
    }

    @GetMapping
    public ResponseEntity<CvGeneratorPreference> get() {
        List<CvGeneratorPreference> preferences = preferenceService.findAll();
        return preferences.isEmpty() ? ResponseEntity.noContent().build() : ResponseEntity.ok(preferences.getFirst());
    }

    @PutMapping
    public ResponseEntity<CvGeneratorPreference> put(@RequestBody CvGeneratorPreference preference) {
        return ResponseEntity.ok(preferenceService.save(preference));
    }
}
