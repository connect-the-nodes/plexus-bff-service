package com.zynchub.digital.hubservice.app.feature.status.service;


import com.zynchub.digital.hubservice.app.feature.status.dto.FeatureResponseDto;
import org.springframework.stereotype.Service;

@Service
public class FeatureServiceImpl  implements FeatureService {

    @Override
    public FeatureResponseDto getStatus() {
        return new FeatureResponseDto("OK", "Service is running");
    }
}
