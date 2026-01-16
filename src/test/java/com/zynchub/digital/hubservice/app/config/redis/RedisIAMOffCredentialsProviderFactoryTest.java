package com.zynchub.digital.hubservice.app.config.redis;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import reactor.test.StepVerifier;

class RedisIAMOffCredentialsProviderFactoryTest {

  @Test
  void shouldResolveCredentials() {
    var redisIAMOffCredentialsProviderFactory = new RedisIAMOffCredentialsProviderFactory();

    var redisStandaloneConfiguration = new RedisStandaloneConfiguration();
    redisStandaloneConfiguration.setHostName("localhost");
    redisStandaloneConfiguration.setPort(6380);

    var credentialsProvider =
        redisIAMOffCredentialsProviderFactory.createCredentialsProvider(
            redisStandaloneConfiguration);
    StepVerifier.create(credentialsProvider.resolveCredentials())
        .assertNext(
            credentials -> {
              assertThat(credentials.getUsername()).isNull();
              assertThat(credentials.getPassword()).isNull();
            })
        .verifyComplete();
  }
}
