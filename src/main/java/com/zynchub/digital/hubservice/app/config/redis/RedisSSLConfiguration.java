package com.zynchub.digital.hubservice.app.config.redis;

import io.lettuce.core.RedisCredentialsProvider;
import io.lettuce.core.SslVerifyMode;
import io.lettuce.core.protocol.ProtocolVersion;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.event.ApplicationReadyEvent;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.event.EventListener;
import org.springframework.data.redis.RedisConnectionFailureException;
import org.springframework.data.redis.connection.RedisConfiguration;
import org.springframework.data.redis.connection.RedisPassword;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import org.springframework.data.redis.connection.lettuce.RedisCredentialsProviderFactory;
import org.springframework.data.redis.connection.lettuce.LettuceClientConfiguration;
import org.springframework.data.redis.connection.lettuce.LettuceConnectionFactory;
import org.springframework.session.data.redis.config.ConfigureRedisAction;
import org.springframework.session.data.redis.config.annotation.web.http.EnableRedisHttpSession;

@Configuration
@ConditionalOnProperty(name = "spring.session.store-type", havingValue = "redis", matchIfMissing = true)
@EnableRedisHttpSession
public class RedisSSLConfiguration {

  private static final Logger log = LoggerFactory.getLogger(RedisSSLConfiguration.class);

  @Value("${spring.data.redis.host:localhost}")
  private String host;

  @Value("${spring.data.redis.port:6379}")
  private int port;

  @Value("${spring.data.redis.ssl.enabled:false}")
  private boolean sslEnabled;

  @Value("${spring.data.redis.iam.enabled:false}")
  private boolean iamEnabled;

  @Value("${spring.data.redis.fail-fast:true}")
  private boolean failFast;

  private final RedisIAMOffCredentialsProviderFactory redisIAMOffCredentialsProviderFactory;
  private final RedisStaticCredentialsProviderFactory redisStaticCredentialsProviderFactory;
  private LettuceConnectionFactory lettuceConnectionFactory;

  public RedisSSLConfiguration(
      RedisIAMOffCredentialsProviderFactory redisIAMOffCredentialsProviderFactory,
      RedisStaticCredentialsProviderFactory redisStaticCredentialsProviderFactory) {
    this.redisIAMOffCredentialsProviderFactory = redisIAMOffCredentialsProviderFactory;
    this.redisStaticCredentialsProviderFactory = redisStaticCredentialsProviderFactory;
  }

  @Bean
  public LettuceConnectionFactory redisConnectionFactory(
      RedisIAMCredentialsProviderFactory redisIAMCredentialsProviderFactory) {
    RedisStandaloneConfiguration standalone = new RedisStandaloneConfiguration(host, port);
    standalone.setPassword(RedisPassword.none());

    RedisCredentialsProvider credentialsProvider =
        iamEnabled
            ? redisIAMCredentialsProviderFactory.createCredentialsProvider(standalone)
            : redisIAMOffCredentialsProviderFactory.createCredentialsProvider(standalone);

    var clientOptions =
        io.lettuce.core.ClientOptions.builder().protocolVersion(ProtocolVersion.RESP3).build();
    LettuceClientConfiguration.LettuceClientConfigurationBuilder builder =
        LettuceClientConfiguration.builder().clientOptions(clientOptions);

    if (sslEnabled) {
      builder.useSsl().verifyPeer(SslVerifyMode.FULL).and();
    }

    builder.redisCredentialsProviderFactory(
        new RedisCredentialsProviderFactory() {
          @Override
          public RedisCredentialsProvider createCredentialsProvider(RedisConfiguration configuration) {
            return credentialsProvider;
          }
        });
    LettuceConnectionFactory factory =
        new LettuceConnectionFactory(standalone, builder.build());
    factory.afterPropertiesSet();
    this.lettuceConnectionFactory = factory;
    return factory;
  }

  @Bean
  public static ConfigureRedisAction configureRedisAction() {
    return ConfigureRedisAction.NO_OP;
  }

  @EventListener
  public void startApp(ApplicationReadyEvent event) {
    if (lettuceConnectionFactory == null) {
      return;
    }
    try {
      lettuceConnectionFactory.getConnection().close();
      log.info("Redis connection successful");
    } catch (RedisConnectionFailureException ex) {
      if (failFast) {
        log.error("Redis connection failed", ex);
        event.getApplicationContext().close();
      } else {
        log.warn("Redis connection failed; continuing without Redis", ex);
      }
    }
  }
}
