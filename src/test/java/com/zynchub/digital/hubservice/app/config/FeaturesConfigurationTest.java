package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import com.zynchub.digital.hubservice.app.mapper.FeatureMapper;
import com.zynchub.digital.hubservice.app.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.service.impl.PropertyFeaturesRetrieverImpl;
import com.zynchub.digital.hubservice.app.service.impl.RemoteFeaturesRetrieverImpl;
import org.junit.jupiter.api.Test;
import org.springframework.boot.autoconfigure.AutoConfigurations;
import org.springframework.boot.autoconfigure.web.client.RestClientAutoConfiguration;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;

class FeaturesConfigurationTest {

  private final ApplicationContextRunner contextRunner =
      new ApplicationContextRunner()
          .withConfiguration(
              AutoConfigurations.of(RestClientAutoConfiguration.class))
          .withUserConfiguration(FeaturesConfiguration.class, FeatureMapper.class);

  @Test
  void property_retriever_is_default() {
    contextRunner
        .withPropertyValues("features.remote.enabled=false", "features.file=test-features.yml")
        .run(
            context -> {
              FeaturesRetriever retriever = context.getBean(FeaturesRetriever.class);
              assertThat(retriever).isInstanceOf(PropertyFeaturesRetrieverImpl.class);
            });
  }

  @Test
  void remote_retriever_is_used_when_enabled() {
    contextRunner
        .withPropertyValues("features.remote.enabled=true", "features.remote.url=https://x")
        .run(
            context -> {
              FeaturesRetriever retriever = context.getBean(FeaturesRetriever.class);
              assertThat(retriever).isInstanceOf(RemoteFeaturesRetrieverImpl.class);
            });
  }
}
