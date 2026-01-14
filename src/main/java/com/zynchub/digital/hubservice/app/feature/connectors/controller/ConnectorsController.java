package com.zynchub.digital.hubservice.app.feature.connectors.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.feature.connectors.dto.ConnectorsResponseDto;
import com.zynchub.digital.hubservice.app.feature.connectors.service.ConnectorsService;
import com.zynchub.digital.hubservice.app.feature.status.dto.FeatureResponseDto;
import com.zynchub.digital.hubservice.app.feature.status.service.FeatureService;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(ApiPaths.API_V3 + "/connectors")
public class ConnectorsController {

    private final ConnectorsService connectorsService;

    public ConnectorsController(ConnectorsService connectorsService) {
        this.connectorsService = connectorsService;
    }

    @GetMapping
    public ApiResponse<ConnectorsResponseDto> status() {
        return ApiResponse.ok(connectorsService.getConnectors());
    }
}