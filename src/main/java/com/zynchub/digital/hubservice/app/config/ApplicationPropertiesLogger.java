package com.zynchub.digital.hubservice.app.config;

import java.util.Locale;
import java.util.Map;
import java.util.TreeMap;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.core.env.ConfigurableEnvironment;
import org.springframework.core.env.Environment;
import org.springframework.core.env.MapPropertySource;
import org.springframework.core.env.PropertySource;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "property-logger")
public class ApplicationPropertiesLogger implements ApplicationRunner {

  private static final Logger log = LoggerFactory.getLogger(ApplicationPropertiesLogger.class);
  private static final String REDACTED = "[!!!REDACTED!!!]";

  private final Environment environment;
  private boolean enabled;

  public ApplicationPropertiesLogger(Environment environment) {
    this.environment = environment;
  }

  public void setEnabled(boolean enabled) {
    this.enabled = enabled;
  }

  @Override
  public void run(ApplicationArguments args) {
    if (!enabled) {
      return;
    }

    log.info("===== Application YAML Properties at Startup =====");
    if (!(environment instanceof ConfigurableEnvironment configurableEnvironment)) {
      return;
    }

    for (PropertySource<?> source : configurableEnvironment.getPropertySources()) {
      if (!(source instanceof MapPropertySource mapSource)) {
        continue;
      }
      String name = source.getName().toLowerCase(Locale.ROOT);
      if (!name.contains("application")) {
        continue;
      }

      Map<String, Object> sorted = new TreeMap<>(mapSource.getSource());
      for (Map.Entry<String, Object> entry : sorted.entrySet()) {
        String key = entry.getKey();
        Object value = entry.getValue();
        String printable = isSensitiveKey(key) ? REDACTED : String.valueOf(value);
        log.info("{} = {}", key, printable);
      }
    }
  }

  private boolean isSensitiveKey(String key) {
    String normalized = key.toLowerCase(Locale.ROOT);
    return normalized.contains("secret")
        || normalized.contains("password")
        || normalized.contains("token")
        || normalized.contains("key");
  }
}
