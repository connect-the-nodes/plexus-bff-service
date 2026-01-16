package com.zynchub.digital.hubservice.app.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConfigurationProperties(prefix = "spring.retry")
public class RetryConfiguration {

  private int maxAttempts;
  private long delay;
  private long maxDelay;
  private boolean random;

  public int getMaxAttempts() {
    return maxAttempts;
  }

  public void setMaxAttempts(int maxAttempts) {
    this.maxAttempts = maxAttempts;
  }

  public long getDelay() {
    return delay;
  }

  public void setDelay(long delay) {
    this.delay = delay;
  }

  public long getMaxDelay() {
    return maxDelay;
  }

  public void setMaxDelay(long maxDelay) {
    this.maxDelay = maxDelay;
  }

  public boolean isRandom() {
    return random;
  }

  public void setRandom(boolean random) {
    this.random = random;
  }
}
