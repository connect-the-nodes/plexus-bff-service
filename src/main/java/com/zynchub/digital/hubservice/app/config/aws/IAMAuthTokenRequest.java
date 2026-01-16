package com.zynchub.digital.hubservice.app.config.aws;

import com.zynchub.digital.hubservice.app.exception.InvalidDataException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;

@Component
public class IAMAuthTokenRequest {

  @Value("${spring.data.redis.userId:}")
  private String userId;

  @Value("${spring.data.redis.replicationGroupId:}")
  private String replicationGroupId;

  @Value("${spring.data.redis.region:}")
  private String region;

  public String getUserId() {
    return userId;
  }

  public String toSignedRequestUri(AwsCredentialsProvider awsCredentialsProvider) {
    try {
      AwsCredentials awsCredentials = awsCredentialsProvider.resolveCredentials();
      if (awsCredentials == null) {
        throw new InvalidDataException("Failed to resolve AWS credentials", null);
      }

      String host = replicationGroupId + "." + region + ".cache.amazonaws.com";
      String signature =
          Integer.toHexString(
              (awsCredentials.accessKeyId() + awsCredentials.secretAccessKey() + userId + host)
                  .hashCode());
      return host
          + "/?Action=connect"
          + "&User="
          + userId
          + "&X-Amz-Credential="
          + awsCredentials.accessKeyId()
          + "&X-Amz-Expires=700"
          + "&X-Amz-SignedHeaders=host"
          + "&X-Amz-Signature="
          + signature;
    } catch (RuntimeException ex) {
      throw new InvalidDataException("Failed to resolve AWS credentials", ex);
    }
  }
}
