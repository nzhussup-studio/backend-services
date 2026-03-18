package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.WorkExperience;
import com.nzhussup.baseservice.repository.WorkExperienceRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class WorkExperienceService extends BaseService<WorkExperience> {


    @Autowired
    public WorkExperienceService(WorkExperienceRepository workExperienceRepository) {
        super(workExperienceRepository);
    }
}
