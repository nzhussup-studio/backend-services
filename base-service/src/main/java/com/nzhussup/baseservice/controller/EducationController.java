package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.Education;
import com.nzhussup.baseservice.service.EducationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(AppConfig.baseApiPath+"education")
public class EducationController extends BaseController<Education>{

    @Autowired
    public EducationController(EducationService educationService) {
        super(educationService);
    }
}
