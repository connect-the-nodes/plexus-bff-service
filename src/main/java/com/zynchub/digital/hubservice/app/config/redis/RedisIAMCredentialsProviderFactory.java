package com.zynchub.digital.hubservice.app.config.redis;

import io.lettuce.core.RedisCredentialsProvider;
import org.springframework.data.redis.connection.RedisStandaloneConfiguration;
import org.springframework.stereotype.Component;

@Component
public class RedisIAMCredentialsProviderFactory {

  private final RedisIAMAuthCredentialsProvider credentialsProvider;

  public RedisIAMCredentialsProviderFactory(RedisIAMAuthCredentialsProvider credentialsProvider) {
    this.credentialsProvider = credentialsProvider;
  }

  public RedisCredentialsProvider createCredentialsProvider(
      RedisStandaloneConfiguration redisStandaloneConfiguration) {
    return credentialsProvider;
  }

  public RedisIAMAuthCredentialsProvider getCredentialsProvider() {
    return credentialsProvider;
  }
}
