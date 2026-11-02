package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.model.*;
import com.nzhussup.baseservice.service.*;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
public class PublicContentController {
    private final WorkExperienceService workExperiences;
    private final SkillService skills;
    private final ProjectService projects;
    private final EducationService educationRecords;
    private final CertificateService certificates;

    public PublicContentController(
            WorkExperienceService workExperiences,
            SkillService skills,
            ProjectService projects,
            EducationService educationRecords,
            CertificateService certificates) {
        this.workExperiences = workExperiences;
        this.skills = skills;
        this.projects = projects;
        this.educationRecords = educationRecords;
        this.certificates = certificates;
    }

    @GetMapping("/v1/work-experiences")
    public ResponseEntity<List<WorkExperience>> workExperiences() {
        return ResponseEntity.ok(workExperiences.findAll());
    }

    @GetMapping("/v1/skills")
    public ResponseEntity<List<Skill>> skills() {
        return ResponseEntity.ok(skills.findAll());
    }

    @GetMapping("/v1/projects")
    public ResponseEntity<List<Project>> projects() {
        return ResponseEntity.ok(projects.findAll());
    }

    @GetMapping("/v1/education-records")
    public ResponseEntity<List<Education>> educationRecords() {
        return ResponseEntity.ok(educationRecords.findAll());
    }

    @GetMapping("/v1/certificates")
    public ResponseEntity<List<Certificate>> certificates() {
        return ResponseEntity.ok(certificates.findAll());
    }
}
