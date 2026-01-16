package com.zynchub.digital.hubservice.app.feature.flags.dto;

import com.zynchub.digital.hubservice.app.feature.flags.model.FeatureFlag;
import java.util.List;

public record FeatureFlagsResponseDto(List<FeatureFlag> features) {}
