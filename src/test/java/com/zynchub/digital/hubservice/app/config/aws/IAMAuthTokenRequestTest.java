package com.zynchub.digital.hubservice.app.config.aws;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.zynchub.digital.hubservice.app.exception.InvalidDataException;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.test.util.ReflectionTestUtils;
import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;

@ExtendWith(MockitoExtension.class)
class IAMAuthTokenRequestTest {

  @Mock private AwsCredentials awsCredentials;
  @Mock private AwsCredentialsProvider awsCredentialsProvider;
  @InjectMocks private IAMAuthTokenRequest iamAuthTokenRequest;

  @Test
  void shouldReturnSignedRequestUri() {
    ReflectionTestUtils.setField(iamAuthTokenRequest, "userId", "testUser");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "replicationGroupId", "testReplicationGroup");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "region", "testRegion");

    when(awsCredentialsProvider.resolveCredentials()).thenReturn(awsCredentials);
    when(awsCredentials.accessKeyId()).thenReturn("accessKey");
    when(awsCredentials.secretAccessKey()).thenReturn("secretKey");
    String signedRequestUri = iamAuthTokenRequest.toSignedRequestUri(awsCredentialsProvider);

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

    when(awsCredentialsProvider.resolveCredentials()).thenThrow(RuntimeException.class);

    assertThrows(
        InvalidDataException.class,
        () -> iamAuthTokenRequest.toSignedRequestUri(awsCredentialsProvider));
  }

  @Test
  void shouldThrowInvalidDataExceptionWhenInterruptedExceptionOccurs() {
    ReflectionTestUtils.setField(iamAuthTokenRequest, "userId", "testUser");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "replicationGroupId", "testReplicationGroup");
    ReflectionTestUtils.setField(iamAuthTokenRequest, "region", "testRegion");

    when(awsCredentialsProvider.resolveCredentials()).thenReturn(null);

    assertThrows(
        InvalidDataException.class,
        () -> iamAuthTokenRequest.toSignedRequestUri(awsCredentialsProvider));
  }

}
