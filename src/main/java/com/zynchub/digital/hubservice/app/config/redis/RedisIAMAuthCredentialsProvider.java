package com.zynchub.digital.hubservice.app.config.redis;

import com.zynchub.digital.hubservice.app.config.aws.IAMAuthTokenRequest;
import io.lettuce.core.RedisCredentials;
import io.lettuce.core.RedisCredentialsProvider;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;

@Component
public class RedisIAMAuthCredentialsProvider implements RedisCredentialsProvider {

  private final IAMAuthTokenRequest iamAuthTokenRequest;
  private final AwsCredentialsProvider awsCredentialsProvider;

  public RedisIAMAuthCredentialsProvider(
      IAMAuthTokenRequest iamAuthTokenRequest, AwsCredentialsProvider awsCredentialsProvider) {
    this.iamAuthTokenRequest = iamAuthTokenRequest;
    this.awsCredentialsProvider = awsCredentialsProvider;
  }

  @Override
  public Mono<RedisCredentials> resolveCredentials() {
    return Mono.fromSupplier(
        () ->
            RedisCredentials.just(
                iamAuthTokenRequest.getUserId(),
                iamAuthTokenRequest.toSignedRequestUri(awsCredentialsProvider)));
  }
}
