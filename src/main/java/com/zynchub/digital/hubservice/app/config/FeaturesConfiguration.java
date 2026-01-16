package com.zynchub.digital.hubservice.app.config;

import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.service.impl.AppConfigFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.service.impl.PropertyFeaturesRetrieverImpl;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.env.Environment;

@Configuration
@RequiredArgsConstructor
public class FeaturesConfiguration {

  private final Environment environment;
  private final FeatureMapper featureMapper;

  @Bean
  @ConditionalOnProperty(
      value = "aws.app-config.features.enabled",
      havingValue = "true",
      matchIfMissing = true)
  public FeaturesRetriever appConfigFeaturesRetriever() {
    return new AppConfigFeaturesRetrieverImpl(environment, featureMapper);
  }

  @Bean
  @ConditionalOnProperty(
      value = "aws.app-config.features.enabled",
      havingValue = "false",
      matchIfMissing = false)
  public FeaturesRetriever propertyFeaturesRetriever() {
    return new PropertyFeaturesRetrieverImpl("features.yml", featureMapper);
  }
}
