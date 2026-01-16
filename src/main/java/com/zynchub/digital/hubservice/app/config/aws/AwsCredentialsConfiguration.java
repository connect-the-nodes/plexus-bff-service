package com.zynchub.digital.hubservice.app.config.aws;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;

@Configuration
public class AwsCredentialsConfiguration {

  @Bean
  public AwsCredentialsProvider awsCredentialsProvider() {
    return DefaultCredentialsProvider.create();
  }
}
