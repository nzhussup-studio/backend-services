package com.nzhussup.baseservice.controller;

import com.nzhussup.baseservice.config.AppConfig;
import com.nzhussup.baseservice.model.Certificate;
import com.nzhussup.baseservice.service.CertificateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping(AppConfig.baseApiPath+"certificate")
public class CertificateController extends BaseController<Certificate> {

    @Autowired
    public CertificateController(CertificateService certificateService) {
        super(certificateService);
    }
}
