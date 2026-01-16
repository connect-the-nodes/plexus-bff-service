package com.zynchub.digital.hubservice.app.feature.flags.controller;

import com.zynchub.digital.hubservice.app.common.ApiPaths;
import com.zynchub.digital.hubservice.app.common.response.ApiResponse;
import com.zynchub.digital.hubservice.app.feature.flags.dto.FeatureFlagsResponseDto;
import com.zynchub.digital.hubservice.app.feature.flags.service.FeatureFlagsService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
@RequestMapping(ApiPaths.API_V1 + "/features")
public class FeatureFlagsController {

  private final FeatureFlagsService featureFlagsService;

  @GetMapping
  public ApiResponse<FeatureFlagsResponseDto> retrieveFeatures() {
    return ApiResponse.ok(featureFlagsService.retrieveFeatures());
  }
}
