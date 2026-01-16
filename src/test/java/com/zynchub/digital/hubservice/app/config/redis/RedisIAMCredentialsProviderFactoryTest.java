package com.zynchub.digital.hubservice.app.config.redis;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

@ExtendWith(MockitoExtension.class)
class RedisIAMCredentialsProviderFactoryTest {

  @InjectMocks private RedisIAMCredentialsProviderFactory redisIAMCredentialsProviderFactory;

  @Mock private RedisIAMAuthCredentialsProvider redisIAMAuthCredentialsProvider;

  @Test
  void shouldReturnInjectedCredentialProvider() {
    assertThat(redisIAMCredentialsProviderFactory.createCredentialsProvider(null))
        .isEqualTo(redisIAMAuthCredentialsProvider);
  }

  @Test
  void shouldReturnInjectedCredentialProviderThroughGetter() {
    assertThat(redisIAMCredentialsProviderFactory.getCredentialsProvider())
        .isEqualTo(redisIAMAuthCredentialsProvider);
  }
}
