package com.zynchub.digital.hubservice.app.feature;

import com.zynchub.digital.hubservice.app.feature.FeatureAssociation;

public class FeatureAspectTestService {

  @FeatureAssociation(name = "TEST_ENABLED")
  public String enabledFeature() {
    return "TEST_ENABLED";
  }

  @FeatureAssociation(name = "TEST_DISABLED")
  public String disabledFeature() {
    return "TEST_DISABLED";
  }

  @FeatureAssociation(name = "TEST_ENABLED", invert = true)
  public String invertedEnabledFeature() {
    return "TEST_ENABLED_INVERTED";
  }

  @FeatureAssociation(name = "TEST_DISABLED", invert = true)
  public String invertedDisabledFeature() {
    return "TEST_DISABLED_INVERTED";
  }

  @FeatureAssociation(name = "TEST_MISSING")
  public String missingFeature() {
    return "TEST_MISSING";
  }
}
