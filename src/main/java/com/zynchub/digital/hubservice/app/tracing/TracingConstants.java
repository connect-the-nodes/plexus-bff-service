package com.zynchub.digital.hubservice.app.tracing;

public final class TracingConstants {

  public static final String CORRELATION_ID_HEADER = "X-Correlation-Id";
  public static final String CORRELATION_ID_BAGGAGE_KEY = "zynchubCorrelationId";

  private TracingConstants() {
    throw new UnsupportedOperationException("Utility class");
  }
}
