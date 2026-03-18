package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.Education;
import com.nzhussup.baseservice.repository.EducationRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class EducationService extends BaseService<Education> {


    @Autowired
    public EducationService(EducationRepository educationRepository) {
        super(educationRepository);
    }
}

