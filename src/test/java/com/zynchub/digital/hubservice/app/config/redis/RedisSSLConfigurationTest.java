package com.zynchub.digital.hubservice.app.config.redis;

import static org.assertj.core.api.SoftAssertions.assertSoftly;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import io.lettuce.core.RedisCredentialsProvider;
import io.lettuce.core.protocol.ProtocolVersion;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.data.redis.RedisConnectionFailureException;
import org.springframework.data.redis.connection.RedisConnection;
import org.springframework.data.redis.connection.RedisPassword;
import org.springframework.data.redis.connection.lettuce.LettuceConnectionFactory;
import org.springframework.session.data.redis.config.ConfigureRedisAction;
import org.springframework.test.util.ReflectionTestUtils;

@ExtendWith(MockitoExtension.class)
class RedisSSLConfigurationTest {

  @Mock private RedisIAMCredentialsProviderFactory redisIAMCredentialsProviderFactory;
  @Mock private RedisIAMOffCredentialsProviderFactory redisIAMOffCredentialsProviderFactory;
  @Mock private RedisStaticCredentialsProviderFactory redisStaticCredentialsProviderFactory;
  @Mock private RedisCredentialsProvider redisCredentialsProvider;
  @Mock private RedisConnection redisConnection;
  @InjectMocks private RedisSSLConfiguration redisSSLConfiguration;

  @Test
  void testConnectionFactoryConfigurationWithSSlEnabled() {
    when(redisIAMCredentialsProviderFactory.createCredentialsProvider(any()))
        .thenReturn(redisCredentialsProvider);
    ReflectionTestUtils.setField(redisSSLConfiguration, "host", "localhost");
    ReflectionTestUtils.setField(redisSSLConfiguration, "port", 6379);
    ReflectionTestUtils.setField(redisSSLConfiguration, "sslEnabled", true);
    ReflectionTestUtils.setField(redisSSLConfiguration, "iamEnabled", true);

    var lettuceConnectionFactory =
        (LettuceConnectionFactory)
            redisSSLConfiguration.redisConnectionFactory(redisIAMCredentialsProviderFactory);

    assertSoftly(
        softly -> {
          var clientConfiguration = lettuceConnectionFactory.getClientConfiguration();
          var clientOptions = clientConfiguration.getClientOptions().get();
          var config = lettuceConnectionFactory.getStandaloneConfiguration();

          softly.assertThat(lettuceConnectionFactory.isUseSsl()).isTrue();
          softly.assertThat(clientOptions.getProtocolVersion()).isEqualTo(ProtocolVersion.RESP3);
          softly.assertThat(config.getPassword()).isEqualTo(RedisPassword.none());
          softly.assertThat(lettuceConnectionFactory.isVerifyPeer()).isTrue();
        });
  }

  @Test
  void testConnectionFactoryConfigurationWithSSlDisabled() {
    when(redisIAMCredentialsProviderFactory.createCredentialsProvider(any()))
        .thenReturn(redisCredentialsProvider);

    ReflectionTestUtils.setField(redisSSLConfiguration, "host", "localhost");
    ReflectionTestUtils.setField(redisSSLConfiguration, "port", 6379);
    ReflectionTestUtils.setField(redisSSLConfiguration, "sslEnabled", false);
    ReflectionTestUtils.setField(redisSSLConfiguration, "iamEnabled", true);

    var lettuceConnectionFactory =
        (LettuceConnectionFactory)
            redisSSLConfiguration.redisConnectionFactory(redisIAMCredentialsProviderFactory);

    assertSoftly(
        softly -> {
          var clientConfiguration = lettuceConnectionFactory.getClientConfiguration();
          var clientOptions = clientConfiguration.getClientOptions().get();
          var config = lettuceConnectionFactory.getStandaloneConfiguration();

          softly.assertThat(lettuceConnectionFactory.isUseSsl()).isFalse();
          softly.assertThat(clientOptions.getProtocolVersion()).isEqualTo(ProtocolVersion.RESP3);
          softly.assertThat(config.getPassword()).isEqualTo(RedisPassword.none());
        });
  }

  @Test
  void configureRedisActionShouldReturnNoOp() {
    ConfigureRedisAction result = RedisSSLConfiguration.configureRedisAction();
    assertEquals(ConfigureRedisAction.NO_OP, result);
  }

  @Test
  void startAppShouldLogConnectionSuccess() {
    var lettuceConnectionFactory = mock(LettuceConnectionFactory.class);
    ReflectionTestUtils.setField(
        redisSSLConfiguration, "lettuceConnectionFactory", lettuceConnectionFactory);
    ReflectionTestUtils.setField(redisSSLConfiguration, "failFast", true);
    when(lettuceConnectionFactory.getConnection()).thenReturn(redisConnection);

    ApplicationReadyEvent applicationReadyEvent = mock(ApplicationReadyEvent.class);
    redisSSLConfiguration.startApp(applicationReadyEvent);

    verify(redisConnection, times(1)).close();
  }

  @Test
  void startAppShouldLogConnectionFailure() {
    var lettuceConnectionFactory = mock(LettuceConnectionFactory.class);
    ReflectionTestUtils.setField(
        redisSSLConfiguration, "lettuceConnectionFactory", lettuceConnectionFactory);
    ReflectionTestUtils.setField(redisSSLConfiguration, "failFast", true);
    when(lettuceConnectionFactory.getConnection()).thenThrow(RedisConnectionFailureException.class);

    ApplicationReadyEvent applicationReadyEvent = mock(ApplicationReadyEvent.class);
    ConfigurableApplicationContext context = mock(ConfigurableApplicationContext.class);
    when(applicationReadyEvent.getApplicationContext()).thenReturn(context);

    redisSSLConfiguration.startApp(applicationReadyEvent);

    verify(context, times(1)).close();
  }
}
