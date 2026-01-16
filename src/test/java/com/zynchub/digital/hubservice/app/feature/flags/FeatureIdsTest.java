package com.zynchub.digital.hubservice.app.feature.flags;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.lang.reflect.Constructor;
import org.junit.jupiter.api.Test;

class FeatureIdsTest {

  @Test
  void constants_are_defined() {
    assertThat(FeatureIds.CONNECTORS).isEqualTo("FEATURE_CONNECTORS");
    assertThat(FeatureIds.CONNECTORS_V2).isEqualTo("FEATURE_CONNECTORS_V2");
  }

  @Test
  void constructor_is_not_accessible() throws Exception {
    Constructor<FeatureIds> constructor = FeatureIds.class.getDeclaredConstructor();
    constructor.setAccessible(true);

    assertThatThrownBy(constructor::newInstance)
        .hasRootCauseInstanceOf(UnsupportedOperationException.class);
  }
}
