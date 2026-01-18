package com.zynchub.digital.hubservice.app.config.aws;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertThrows;
import com.zynchub.digital.hubservice.app.exception.InvalidDataException;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import org.junit.jupiter.api.Test;
import org.springframework.test.util.ReflectionTestUtils;
import software.amazon.awssdk.auth.credentials.AwsBasicCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.auth.credentials.StaticCredentialsProvider;
import software.amazon.awssdk.identity.spi.AwsCredentialsIdentity;

class IAMAuthTokenRequestTest {

  private final IAMAuthTokenRequest iamAuthTokenRequest = new IAMAuthTokenRequest();

  @Test
  void shouldReturnSignedRequestUri() {
    ReflectionTestUtils.setField(iamAuthTokenRequest, "userId", "testUser");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "replicationGroupId", "testReplicationGroup");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "region", "testRegion");

    AwsCredentialsProvider provider =
        StaticCredentialsProvider.create(AwsBasicCredentials.create("accessKey", "secretKey"));
    String signedRequestUri = iamAuthTokenRequest.toSignedRequestUri(provider);

    assertAll(
        () -> assertThat(signedRequestUri).doesNotContain("https://"),
        () -> assertThat(signedRequestUri).contains("testReplicationGroup"),
        () -> assertThat(signedRequestUri).contains("Action=connect"),
        () -> assertThat(signedRequestUri).contains("User=testUser"),
        () -> assertThat(signedRequestUri).contains("X-Amz-Credential=accessKey"),
        () -> assertThat(signedRequestUri).contains("X-Amz-Expires=700"),
        () -> assertThat(signedRequestUri).contains("X-Amz-SignedHeaders=host"),
        () -> assertThat(signedRequestUri).contains("X-Amz-Signature="));
  }

  @Test
  void shouldThrowInvalidDataExceptionWhenExecutionExceptionOccurs() {
    ReflectionTestUtils.setField(iamAuthTokenRequest, "userId", "testUser");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "replicationGroupId", "testReplicationGroup");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "region", "testRegion");

    assertThrows(
        InvalidDataException.class,
        () ->
            iamAuthTokenRequest.toSignedRequestUri(
                failingIdentityProvider(new ExecutionException(new RuntimeException("boom")))));
  }

  @Test
  void shouldThrowInvalidDataExceptionWhenInterruptedExceptionOccurs() {
    ReflectionTestUtils.setField(iamAuthTokenRequest, "userId", "testUser");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "replicationGroupId", "testReplicationGroup");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "region", "testRegion");

    assertThrows(
        InvalidDataException.class,
        () ->
            iamAuthTokenRequest.toSignedRequestUri(
                failingIdentityProvider(new InterruptedException("interrupted"))));
  }

  private static AwsCredentialsProvider failingIdentityProvider(Throwable throwable) {
    CompletableFuture<AwsCredentialsIdentity> failed = new CompletableFuture<>();
    failed.completeExceptionally(throwable);
    return new AwsCredentialsProvider() {
      @Override
      public AwsCredentials resolveCredentials() {
        return AwsBasicCredentials.create("accessKey", "secretKey");
      }

      @Override
      public CompletableFuture<? extends AwsCredentialsIdentity> resolveIdentity() {
        return failed;
      }
    };
  }
}
