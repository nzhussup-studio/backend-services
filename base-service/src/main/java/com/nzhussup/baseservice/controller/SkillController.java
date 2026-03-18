package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.Skill;
import com.nzhussup.baseservice.service.SkillService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(AppConfig.baseApiPath+"skill")
public class SkillController extends BaseController<Skill> {

    @Autowired
    public SkillController(SkillService skillService) {
        super(skillService);
    }
}
