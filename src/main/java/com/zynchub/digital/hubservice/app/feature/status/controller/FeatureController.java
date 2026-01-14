package com.zynchub.digital.hubservice.app.feature.status.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.feature.status.dto.FeatureResponseDto;
import com.zynchub.digital.hubservice.app.feature.status.service.FeatureService;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(ApiPaths.API_V2 + "/status")
public class FeatureController {

    private final FeatureService statusService;

    public FeatureController(FeatureService statusService) {
        this.statusService = statusService;
    }

    @GetMapping
    public ApiResponse<FeatureResponseDto> status() {
        return ApiResponse.ok(statusService.getStatus());
    }
}
