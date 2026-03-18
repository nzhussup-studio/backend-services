package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.Project;
import com.nzhussup.baseservice.repository.ProjectRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;


@Service
public class ProjectService extends BaseService<Project> {

    @Autowired
    public ProjectService(ProjectRepository projectRepository) {
        super(projectRepository);
    }

}
