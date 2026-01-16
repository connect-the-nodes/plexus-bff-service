package com.zynchub.digital.hubservice.app.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.dto.ConnectorsResponseDto;
import com.zynchub.digital.hubservice.app.service.ConnectorsService;
import com.zynchub.digital.hubservice.app.feature.FeatureIds;
import com.zynchub.digital.hubservice.app.feature.FeatureAssociation;
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
    @FeatureAssociation(name = FeatureIds.CONNECTORS)
    public ApiResponse<ConnectorsResponseDto> status() {
        return ApiResponse.ok(connectorsService.getConnectors());
    }
}
