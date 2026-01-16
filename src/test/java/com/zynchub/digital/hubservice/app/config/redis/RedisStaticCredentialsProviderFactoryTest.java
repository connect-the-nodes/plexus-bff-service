package com.zynchub.digital.hubservice.app.config.redis;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import reactor.test.StepVerifier;

class RedisStaticCredentialsProviderFactoryTest {

  @Test
  void shouldResolveCredentials() {
    var redisStaticCredentialsProviderFactory = new RedisStaticCredentialsProviderFactory();

    var redisStandaloneConfiguration = new RedisStandaloneConfiguration();
    redisStandaloneConfiguration.setHostName("localhost");
    redisStandaloneConfiguration.setPort(6380);

    var credentialsProvider =
        redisStaticCredentialsProviderFactory.createCredentialsProvider(
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
