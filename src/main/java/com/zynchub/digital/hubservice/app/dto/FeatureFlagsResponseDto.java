package com.zynchub.digital.hubservice.app.dto;

import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import java.util.List;

public record FeatureFlagsResponseDto(List<FeatureFlag> features) {}
