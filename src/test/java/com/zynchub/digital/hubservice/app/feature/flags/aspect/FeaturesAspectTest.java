package com.zynchub.digital.hubservice.app.feature.flags.aspect;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import com.zynchub.digital.hubservice.app.feature.flags.exception.FeatureNotFoundException;
import com.zynchub.digital.hubservice.app.feature.flags.service.FeaturesRetriever;
import com.zynchub.digital.hubservice.app.feature.flags.service.impl.PropertyFeaturesRetrieverImpl;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.security.access.AccessDeniedException;

@SpringBootTest(
    properties = {
      "security.enabled=false",
      "features.remote.enabled=false",
      "features.file=test-features.yml"
    })
class FeaturesAspectTest {

  @Autowired private FeatureAspectTestService featureAspectTestService;

  @Test
  void allows_enabled_feature() {
    assertThat(featureAspectTestService.enabledFeature()).isEqualTo("TEST_ENABLED");
  }

  @Test
  void denies_disabled_feature() {
    assertThatThrownBy(featureAspectTestService::disabledFeature)
        .isInstanceOf(AccessDeniedException.class)
        .hasMessage("Feature TEST_DISABLED is not enabled!");
  }

  @Test
  void inverted_enabled_feature_is_denied() {
    assertThatThrownBy(featureAspectTestService::invertedEnabledFeature)
        .isInstanceOf(AccessDeniedException.class)
        .hasMessage("Feature TEST_ENABLED is not enabled!");
  }

  @Test
  void inverted_disabled_feature_is_allowed() {
    assertThat(featureAspectTestService.invertedDisabledFeature())
        .isEqualTo("TEST_DISABLED_INVERTED");
  }

  @Test
  void missing_feature_throws_not_found() {
    assertThatThrownBy(featureAspectTestService::missingFeature)
        .isInstanceOf(FeatureNotFoundException.class);
  }

  @TestConfiguration
  static class FeatureAspectTestConfiguration {
    @Bean
    FeatureAspectTestService featureAspectTestService() {
      return new FeatureAspectTestService();
    }

    @Bean
    @org.springframework.context.annotation.Primary
    FeaturesRetriever testFeaturesRetriever(
        com.zynchub.digital.hubservice.app.feature.flags.mapper.FeatureMapper featureMapper) {
      return new PropertyFeaturesRetrieverImpl("test-features.yml", featureMapper);
    }
  }
}
