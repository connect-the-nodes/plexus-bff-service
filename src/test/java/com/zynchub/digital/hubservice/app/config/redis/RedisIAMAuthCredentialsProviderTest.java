package com.zynchub.digital.hubservice.app.config.redis;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.Mockito.when;

import com.zynchub.digital.hubservice.app.config.aws.IAMAuthTokenRequest;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import reactor.test.StepVerifier;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;

@ExtendWith(MockitoExtension.class)
class RedisIAMAuthCredentialsProviderTest {

  @InjectMocks private RedisIAMAuthCredentialsProvider redisIAMAuthCredentialsProvider;

  @Mock private IAMAuthTokenRequest iamAuthTokenRequest;

  @Mock private AwsCredentialsProvider awsCredentialsProvider;

  @Test
  void shouldResolveCredentials() {
    when(iamAuthTokenRequest.toSignedRequestUri(awsCredentialsProvider)).thenReturn("token");
    when(iamAuthTokenRequest.getUserId()).thenReturn("testUser");
    StepVerifier.create(redisIAMAuthCredentialsProvider.resolveCredentials())
        .assertNext(
            credentials -> {
              assertThat(credentials.getUsername()).isEqualTo("testUser");
              assertThat(new String(credentials.getPassword())).isEqualTo("token");
            })
        .verifyComplete();
  }
}
