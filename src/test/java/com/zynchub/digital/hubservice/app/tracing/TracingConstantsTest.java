package com.zynchub.digital.hubservice.app.tracing;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.lang.reflect.Constructor;
import org.junit.jupiter.api.Test;

class TracingConstantsTest {

  @Test
  void constants_are_defined() {
    assertThat(TracingConstants.CORRELATION_ID_HEADER).isEqualTo("X-Correlation-Id");
    assertThat(TracingConstants.CORRELATION_ID_BAGGAGE_KEY).isEqualTo("zynchubCorrelationId");
  }

  @Test
  void constructor_is_not_accessible() throws Exception {
    Constructor<TracingConstants> constructor = TracingConstants.class.getDeclaredConstructor();
    constructor.setAccessible(true);

    assertThatThrownBy(constructor::newInstance)
        .hasRootCauseInstanceOf(UnsupportedOperationException.class);
  }
}
