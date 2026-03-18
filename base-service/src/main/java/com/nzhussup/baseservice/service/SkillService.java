package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.Skill;
import com.nzhussup.baseservice.repository.SkillRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class SkillService extends BaseService<Skill> {

    @Autowired
    public SkillService(SkillRepository skillRepository) {
        super(skillRepository);
    }
}

