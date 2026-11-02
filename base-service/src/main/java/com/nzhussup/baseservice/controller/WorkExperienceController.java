package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.WorkExperience;
import com.nzhussup.baseservice.service.WorkExperienceService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(AppConfig.adminApiPath + "work-experiences")
public class WorkExperienceController extends BaseController<WorkExperience> {

    @Autowired
    public WorkExperienceController(WorkExperienceService workExperienceService) {
        super(workExperienceService);
    }
}
