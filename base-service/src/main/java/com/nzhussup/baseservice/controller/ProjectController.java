package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.Project;
import com.nzhussup.baseservice.service.ProjectService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(AppConfig.baseApiPath+"project")
public class ProjectController extends BaseController<Project> {

    @Autowired
    public ProjectController(ProjectService projectService) {
        super(projectService);
    }
}
