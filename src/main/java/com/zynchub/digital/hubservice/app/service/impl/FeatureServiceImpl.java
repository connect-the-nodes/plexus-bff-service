package com.zynchub.digital.hubservice.app.service.impl;

import com.zynchub.digital.hubservice.app.dto.FeatureResponseDto;
import com.zynchub.digital.hubservice.app.service.FeatureService;
import org.springframework.stereotype.Service;

@Service
public class FeatureServiceImpl implements FeatureService {

    @Override
    public FeatureResponseDto getStatus() {
        return new FeatureResponseDto("OK", "Service is running");
    }
}
