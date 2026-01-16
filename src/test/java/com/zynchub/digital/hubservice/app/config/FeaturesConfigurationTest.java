package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.service.impl.AppConfigFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.service.impl.PropertyFeaturesRetrieverImpl;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;

class FeaturesConfigurationTest {

  private final ApplicationContextRunner contextRunner =
      new ApplicationContextRunner()
          .withUserConfiguration(FeaturesConfiguration.class, FeatureMapper.class);

  @Test
  void appconfig_retriever_is_default() {
    contextRunner
        .withPropertyValues("aws.app-config.features.enabled=true")
        .run(
            context -> {
              FeaturesRetriever retriever = context.getBean(FeaturesRetriever.class);
              assertThat(retriever).isInstanceOf(AppConfigFeaturesRetrieverImpl.class);
            });
  }

  @Test
  void property_retriever_is_used_when_disabled() {
    contextRunner
        .withPropertyValues("aws.app-config.features.enabled=false")
        .run(
            context -> {
              FeaturesRetriever retriever = context.getBean(FeaturesRetriever.class);
              assertThat(retriever).isInstanceOf(PropertyFeaturesRetrieverImpl.class);
            });
  }
}
