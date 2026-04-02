package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.CvGeneratorPreference;
import com.nzhussup.baseservice.service.CvGeneratorPreferenceService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(AppConfig.baseApiPath + "cv-generator-preference")
public class CvGeneratorPreferenceController extends BaseController<CvGeneratorPreference> {

    @Autowired
    public CvGeneratorPreferenceController(CvGeneratorPreferenceService preferenceService) {
        super(preferenceService);
    }
}
