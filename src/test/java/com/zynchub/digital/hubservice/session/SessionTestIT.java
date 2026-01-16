package com.zynchub.digital.hubservice.session;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;

import jakarta.servlet.http.HttpSession;
import java.util.Set;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.connection.RedisConnection;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@SpringBootTest
@AutoConfigureMockMvc
@ActiveProfiles("local-it")
class SessionTestIT {

  @Autowired private MockMvc mockMvc;
  @Autowired private RedisConnectionFactory redisConnectionFactory;

  @BeforeEach
  @AfterEach
  void flushKeys() {
    try (RedisConnection connection = redisConnectionFactory.getConnection()) {
      connection.serverCommands().flushAll();
    }
  }

  @Test
  @WithMockUser(username = "test", roles = "USER")
  void sessionIsPersisted() throws Exception {
    assertThat(fetchKeys()).isEmpty();

    mockMvc.perform(post("/_test/session"));

    assertThat(fetchKeys()).isNotEmpty();
  }

  private Set<byte[]> fetchKeys() {
    try (RedisConnection connection = redisConnectionFactory.getConnection()) {
      return connection.keyCommands().keys("*".getBytes());
    }
  }

  @Configuration
  static class TestControllerConfig {
    @Bean
    TestSessionController testSessionController() {
      return new TestSessionController();
    }
  }

  @RestController
  static class TestSessionController {
    @PostMapping("/_test/session")
    ResponseEntity<Void> createSession(HttpSession session) {
      session.setAttribute("test-key", "test-value");
      return ResponseEntity.ok().build();
    }
  }
}
