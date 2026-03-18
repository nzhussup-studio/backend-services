package com.nzhussup.baseservice.service;

import com.nzhussup.baseservice.model.Certificate;
import com.nzhussup.baseservice.repository.CertificateRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class CertificateService extends BaseService<Certificate> {

    @Autowired
    public CertificateService(CertificateRepository certificateRepository) {
        super(certificateRepository);
    }
}
