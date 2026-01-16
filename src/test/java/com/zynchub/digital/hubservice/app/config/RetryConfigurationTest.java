package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;
import org.springframework.boot.autoconfigure.context.ConfigurationPropertiesAutoConfiguration;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;

class RetryConfigurationTest {

  private final ApplicationContextRunner contextRunner =
      new ApplicationContextRunner()
          .withUserConfiguration(RetryConfiguration.class)
          .withConfiguration(
              org.springframework.boot.autoconfigure.AutoConfigurations.of(
                  ConfigurationPropertiesAutoConfiguration.class));

  @Test
  void binds_retry_properties() {
    contextRunner
        .withPropertyValues(
            "spring.retry.max-attempts=3",
            "spring.retry.delay=1000",
            "spring.retry.max-delay=5000",
            "spring.retry.random=true")
        .run(
            context -> {
              RetryConfiguration config = context.getBean(RetryConfiguration.class);
              assertThat(config.getMaxAttempts()).isEqualTo(3);
              assertThat(config.getDelay()).isEqualTo(1000L);
              assertThat(config.getMaxDelay()).isEqualTo(5000L);
              assertThat(config.isRandom()).isTrue();
            });
  }
}
