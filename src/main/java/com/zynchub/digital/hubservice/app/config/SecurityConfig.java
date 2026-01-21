package com.zynchub.digital.hubservice.app.config;

import com.zynchub.digital.hubservice.app.security.JwtAuthConverter;
import com.zynchub.digital.hubservice.app.tracing.CorrelationIdFilter;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.context.SecurityContextHolderFilter;

@Configuration
@ConditionalOnProperty(name = "security.enabled", havingValue = "true", matchIfMissing = true)
public class SecurityConfig {

  @Bean
  SecurityFilterChain filterChain(
      HttpSecurity http, JwtAuthConverter jwtAuthConverter, CorrelationIdFilter correlationIdFilter)
      throws Exception {
    http.csrf(AbstractHttpConfigurer::disable)
        .authorizeHttpRequests(
            auth ->
                auth.requestMatchers(
                        "/actuator/health",
                        "/swagger-ui/**",
                        "/v3/api-docs/**",
                        "/api/v1/hello",
                        "/auth/login",
                        "/auth/callback")
                    .permitAll()
                    .anyRequest()
                    .authenticated())
        .oauth2ResourceServer(
            oauth2 -> oauth2.jwt(jwt -> jwt.jwtAuthenticationConverter(jwtAuthConverter)));

    http.addFilterBefore(correlationIdFilter, SecurityContextHolderFilter.class);

    return http.build();
  }
}
