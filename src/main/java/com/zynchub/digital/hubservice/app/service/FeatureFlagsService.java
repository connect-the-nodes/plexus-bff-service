package com.zynchub.digital.hubservice.app.service;

import com.zynchub.digital.hubservice.app.dto.FeatureFlagsResponseDto;
import com.zynchub.digital.hubservice.app.model.FeatureFlag;
import java.util.List;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Slf4j
@Service
@RequiredArgsConstructor
public class FeatureFlagsService {

  private final FeaturesRetriever featuresRetriever;

  public FeatureFlagsResponseDto retrieveFeatures() {
    List<FeatureFlag> features = featuresRetriever.retrieveFeatures();
    log.debug("Features retrieved are {}", features);
    return new FeatureFlagsResponseDto(features);
  }
}
