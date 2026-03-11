package com.zynchub.digital.hubservice.app.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@EnableConfigurationProperties(ObservabilityClientProperties.Properties.class)
public class ObservabilityClientProperties {

  @ConfigurationProperties(prefix = "observability.service")
  public record Properties(String baseUrl) {}
}
