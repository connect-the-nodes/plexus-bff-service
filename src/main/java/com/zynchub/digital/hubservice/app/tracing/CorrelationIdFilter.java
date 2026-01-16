package com.zynchub.digital.hubservice.app.tracing;

import static com.zynchub.digital.hubservice.app.tracing.TracingConstants.CORRELATION_ID_BAGGAGE_KEY;
import static com.zynchub.digital.hubservice.app.tracing.TracingConstants.CORRELATION_ID_HEADER;

import io.micrometer.tracing.BaggageInScope;
import io.micrometer.tracing.Tracer;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.util.Optional;
import java.util.UUID;
import org.slf4j.MDC;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

@Component
public class CorrelationIdFilter extends OncePerRequestFilter {

  private final Optional<Tracer> tracer;

  public CorrelationIdFilter(Optional<Tracer> tracer) {
    this.tracer = tracer;
  }

  @Override
  protected void doFilterInternal(
      HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
      throws ServletException, IOException {
    String correlationId = request.getHeader(CORRELATION_ID_HEADER);
    if (correlationId == null || correlationId.isBlank()) {
      correlationId = UUID.randomUUID().toString();
    }

    BaggageInScope baggageScope = null;
    try {
      MDC.put(CORRELATION_ID_BAGGAGE_KEY, correlationId);
      if (tracer.isPresent()) {
        baggageScope = tracer.get().createBaggageInScope(CORRELATION_ID_BAGGAGE_KEY, correlationId);
      }
      response.setHeader(CORRELATION_ID_HEADER, correlationId);
      filterChain.doFilter(request, response);
    } finally {
      if (baggageScope != null) {
        baggageScope.close();
      }
      MDC.remove(CORRELATION_ID_BAGGAGE_KEY);
    }
  }
}
