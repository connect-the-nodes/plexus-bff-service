package com.zynchub.digital.hubservice.app.security;

import static org.hamcrest.Matchers.not;
import static org.hamcrest.Matchers.nullValue;
import static org.springframework.security.test.web.servlet.request.SecurityMockMvcRequestPostProcessors.jwt;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.header;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Import;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;

@SpringBootTest(
    properties = {
      "spring.session.store-type=none",
      "management.health.redis.enabled=false",
      "spring.autoconfigure.exclude=org.springframework.boot.autoconfigure.session.SessionAutoConfiguration,org.springframework.boot.autoconfigure.data.redis.RedisAutoConfiguration,org.springframework.boot.autoconfigure.data.redis.RedisRepositoriesAutoConfiguration"
    })
@AutoConfigureMockMvc
@ActiveProfiles("local-it")
@Import(TestJwtDecoderConfig.class)
class SecurityAndTracingTest {

  @Autowired private MockMvc mockMvc;

  @Test
  void healthIsPermittedWithoutAuth() throws Exception {
    mockMvc.perform(get("/actuator/health")).andExpect(status().isOk());
  }

  @Test
  void protectedEndpointIsBlockedWithoutAuth() throws Exception {
    mockMvc.perform(get("/api/v2/status")).andExpect(status().isUnauthorized());
  }

  @Test
  void protectedEndpointIsAllowedWithJwt() throws Exception {
    mockMvc
        .perform(
            get("/api/v2/status")
                .with(jwt().jwt(jwt -> jwt.claim("cognito:groups", List.of("ADMIN")))))
        .andExpect(status().isOk());
  }

  @Test
  void correlationIdIsEchoedForRequests() throws Exception {
    mockMvc
        .perform(get("/actuator/health").header("X-Correlation-Id", "test-correlation-id"))
        .andExpect(status().isOk())
        .andExpect(header().string("X-Correlation-Id", "test-correlation-id"));
  }

  @Test
  void correlationIdIsGeneratedWhenMissing() throws Exception {
    mockMvc
        .perform(get("/actuator/health"))
        .andExpect(status().isOk())
        .andExpect(header().string("X-Correlation-Id", not(nullValue())));
  }
}
