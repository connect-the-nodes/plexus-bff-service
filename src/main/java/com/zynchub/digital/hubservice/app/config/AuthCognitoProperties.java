package com.zynchub.digital.hubservice.app.config;

import java.util.List;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties(AuthCognitoProperties.Properties.class)
public class AuthCognitoProperties {

  @ConfigurationProperties(prefix = "auth.cognito")
  public record Properties(
      boolean enabled,
      String domain,
      String clientId,
      String redirectUri,
      String postLoginRedirectUri,
      List<String> scopes) {}
}
