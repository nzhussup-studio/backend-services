package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.Project;
import com.nzhussup.baseservice.repository.ProjectRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;


@Service
public class ProjectService extends BaseService<Project> {

    private final GitHubProjectDescriptionService gitHubDescriptionService;

    @Autowired
    public ProjectService(ProjectRepository projectRepository,
                          GitHubProjectDescriptionService gitHubDescriptionService) {
        super(projectRepository);
        this.gitHubDescriptionService = gitHubDescriptionService;
    }

    @Override
    public List<Project> findAll() {
        List<Project> projects = repository.findAll();
        projects.forEach(this::addGitHubDescription);
        return projects;
    }

    private void addGitHubDescription(Project project) {
        String description = gitHubDescriptionService.getDescription(project.getUrl());
        project.setDescription(description.isBlank() ? null : description);
    }

}
