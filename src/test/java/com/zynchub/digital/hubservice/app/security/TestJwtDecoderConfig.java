package com.zynchub.digital.hubservice.app.security;

import java.time.Instant;
import java.util.Map;
import org.springframework.boot.test.context.TestConfiguration;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Primary;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.security.oauth2.jwt.JwtDecoder;

@TestConfiguration
public class TestJwtDecoderConfig {

  @Bean
  @Primary
  JwtDecoder jwtDecoder() {
    return token ->
        Jwt.withTokenValue(token)
            .header("alg", "none")
            .claims(claims -> claims.putAll(Map.of("sub", "test-user")))
            .issuedAt(Instant.now())
            .expiresAt(Instant.now().plusSeconds(3600))
            .build();
  }
}
