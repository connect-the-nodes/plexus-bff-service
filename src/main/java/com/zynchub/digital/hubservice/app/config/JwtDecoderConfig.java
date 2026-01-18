package com.zynchub.digital.hubservice.app.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.JwtException;
import org.springframework.security.oauth2.jwt.NimbusJwtDecoder;

@Configuration
@ConditionalOnProperty(name = "security.enabled", havingValue = "true", matchIfMissing = true)
public class JwtDecoderConfig {

  @Bean
  @ConditionalOnMissingBean(JwtDecoder.class)
  JwtDecoder jwtDecoder(
      @Value("${security.jwt.jwk-set-uri:}") String jwkSetUri,
      @Value("${security.jwt.issuer-uri:}") String issuerUri) {
    if (jwkSetUri != null && !jwkSetUri.isBlank()) {
      return NimbusJwtDecoder.withJwkSetUri(jwkSetUri).build();
    }
    if (issuerUri != null && !issuerUri.isBlank()) {
      return NimbusJwtDecoder.withIssuerLocation(issuerUri).build();
    }
    return new JwtDecoder() {
      @Override
      public Jwt decode(String token) {
        throw new JwtException("No JwtDecoder configured; set security.jwt.jwk-set-uri or issuer-uri");
      }
    };
  }
}
