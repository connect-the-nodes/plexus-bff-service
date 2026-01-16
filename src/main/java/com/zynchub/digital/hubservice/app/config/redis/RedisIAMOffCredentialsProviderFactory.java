package com.zynchub.digital.hubservice.app.config.redis;

import io.lettuce.core.RedisCredentials;
import io.lettuce.core.RedisCredentialsProvider;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

@Component
public class RedisIAMOffCredentialsProviderFactory {

  public RedisCredentialsProvider createCredentialsProvider(
      RedisStandaloneConfiguration redisStandaloneConfiguration) {
    return new RedisCredentialsProvider() {
      @Override
      public Mono<RedisCredentials> resolveCredentials() {
        return Mono.just(new SimpleRedisCredentials(null, null));
      }
    };
  }
}
