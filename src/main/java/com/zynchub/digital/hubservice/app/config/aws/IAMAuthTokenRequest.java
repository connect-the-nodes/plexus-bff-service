package com.zynchub.digital.hubservice.app.config.aws;

import static software.amazon.awssdk.http.auth.aws.signer.AwsV4FamilyHttpSigner.AUTH_LOCATION;
import static software.amazon.awssdk.http.auth.aws.signer.AwsV4FamilyHttpSigner.EXPIRATION_DURATION;
import static software.amazon.awssdk.http.auth.aws.signer.AwsV4FamilyHttpSigner.SERVICE_SIGNING_NAME;
import static software.amazon.awssdk.http.auth.aws.signer.AwsV4HttpSigner.REGION_NAME;

import com.zynchub.digital.hubservice.app.exception.InvalidDataException;
import java.io.ByteArrayInputStream;
import java.net.URI;
import java.time.Duration;
import java.util.concurrent.ExecutionException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.http.ContentStreamProvider;
import software.amazon.awssdk.http.SdkHttpFullRequest;
import software.amazon.awssdk.http.SdkHttpMethod;
import software.amazon.awssdk.http.SdkHttpRequest;
import software.amazon.awssdk.http.auth.aws.signer.AwsV4FamilyHttpSigner.AuthLocation;
import software.amazon.awssdk.http.auth.aws.signer.AwsV4HttpSigner;

@Component
public class IAMAuthTokenRequest {
  private static final SdkHttpMethod REQUEST_METHOD = SdkHttpMethod.GET;
  private static final String REQUEST_PROTOCOL = "https://";
  private static final String PARAM_ACTION = "Action";
  private static final String PARAM_USER = "User";
  private static final String ACTION_NAME = "connect";
  private static final String SERVICE_NAME = "elasticache";
  private static final Duration TOKEN_EXPIRY_DURATION_SECONDS = Duration.ofSeconds(700);

  @Value("${spring.data.redis.userId:}")
  private String userId;

  @Value("${spring.data.redis.host:}")
  private String host;

  @Value("${spring.data.redis.replicationGroupId:}")
  private String replicationGroupId;

  @Value("${spring.data.redis.region:}")
  private String region;

  public String getUserId() {
    return userId;
  }

  public String toSignedRequestUri(AwsCredentialsProvider awsCredentialsProvider) {
    SdkHttpRequest request = getSignableRequest();
    try {
      request = sign(request, awsCredentialsProvider);
    } catch (InterruptedException e) {
      Thread.currentThread().interrupt();
      throw new InvalidDataException("Error occurred while signing the request", e);
    } catch (ExecutionException e) {
      throw new InvalidDataException("Error occurred while signing the request", e);
    }
    return request.getUri().toString().replace(REQUEST_PROTOCOL, "");
  }

  private SdkHttpFullRequest getSignableRequest() {
    return SdkHttpFullRequest.builder()
        .method(REQUEST_METHOD)
        .uri(getRequestUri())
        .appendRawQueryParameter(PARAM_ACTION, ACTION_NAME)
        .appendRawQueryParameter(PARAM_USER, userId)
        .build();
  }

  private URI getRequestUri() {
    String resolvedHost = replicationGroupId;
    if (resolvedHost == null || resolvedHost.isBlank()) {
      resolvedHost = host;
    }
    return URI.create(String.format("%s%s/", REQUEST_PROTOCOL, resolvedHost));
  }

  private SdkHttpRequest sign(SdkHttpRequest request, AwsCredentialsProvider awsCredentialsProvider)
      throws ExecutionException, InterruptedException {
    var signer = AwsV4HttpSigner.create();
    var identity = awsCredentialsProvider.resolveIdentity().get();
    var signedRequest =
        signer.sign(
            r ->
                r.identity(identity)
                    .request(request)
                    .payload(() -> new ByteArrayInputStream(new byte[0]))
                    .putProperty(REGION_NAME, region)
                    .putProperty(SERVICE_SIGNING_NAME, SERVICE_NAME)
                    .putProperty(AUTH_LOCATION, AuthLocation.QUERY_STRING)
                    .putProperty(EXPIRATION_DURATION, TOKEN_EXPIRY_DURATION_SECONDS));

    return signedRequest.request();
  }
}
