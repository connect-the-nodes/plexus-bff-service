package com.zynchub.digital.hubservice.app.tracing;

import static com.zynchub.digital.hubservice.app.tracing.TracingConstants.CORRELATION_ID_BAGGAGE_KEY;
import static com.zynchub.digital.hubservice.app.tracing.TracingConstants.CORRELATION_ID_HEADER;
import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import io.micrometer.tracing.Tracer;
import io.micrometer.tracing.BaggageInScope;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletRequest;
import jakarta.servlet.ServletResponse;
import java.io.IOException;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicReference;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.slf4j.MDC;
import org.springframework.mock.web.MockHttpServletRequest;
import org.springframework.mock.web.MockHttpServletResponse;

class CorrelationIdFilterTest {

  @AfterEach
  void clearMdc() {
    MDC.clear();
  }

  @Test
  void uses_request_header_and_sets_response_and_baggage() throws Exception {
    Tracer tracer = mock(Tracer.class);
    BaggageInScope scope = mock(BaggageInScope.class);
    when(tracer.createBaggageInScope(CORRELATION_ID_BAGGAGE_KEY, "incoming-id"))
        .thenReturn(scope);

    var filter = new CorrelationIdFilter(Optional.of(tracer));
    var request = new MockHttpServletRequest();
    request.addHeader(CORRELATION_ID_HEADER, "incoming-id");
    var response = new MockHttpServletResponse();
    var seenCorrelation = new AtomicReference<String>();

    FilterChain chain =
        new FilterChain() {
          @Override
          public void doFilter(ServletRequest req, ServletResponse res) throws IOException {
            seenCorrelation.set(MDC.get(CORRELATION_ID_BAGGAGE_KEY));
          }
        };

    filter.doFilter(request, response, chain);

    assertThat(seenCorrelation.get()).isEqualTo("incoming-id");
    assertThat(response.getHeader(CORRELATION_ID_HEADER)).isEqualTo("incoming-id");
    assertThat(MDC.get(CORRELATION_ID_BAGGAGE_KEY)).isNull();
    verify(tracer).createBaggageInScope(CORRELATION_ID_BAGGAGE_KEY, "incoming-id");
    verify(scope).close();
  }

  @Test
  void generates_correlation_id_when_missing() throws Exception {
    var filter = new CorrelationIdFilter(Optional.empty());
    var request = new MockHttpServletRequest();
    var response = new MockHttpServletResponse();
    var seenCorrelation = new AtomicReference<String>();

    filter.doFilter(
        request,
        response,
        (req, res) -> seenCorrelation.set(MDC.get(CORRELATION_ID_BAGGAGE_KEY)));

    assertThat(seenCorrelation.get()).isNotBlank();
    assertThat(response.getHeader(CORRELATION_ID_HEADER)).isEqualTo(seenCorrelation.get());
    assertThat(MDC.get(CORRELATION_ID_BAGGAGE_KEY)).isNull();
  }
}
