package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.HashMap;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.boot.DefaultApplicationArguments;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;
import org.springframework.core.env.MapPropertySource;
import org.springframework.core.env.StandardEnvironment;
import org.springframework.test.util.ReflectionTestUtils;

@ExtendWith(OutputCaptureExtension.class)
class ApplicationPropertiesLoggerTest {

  @Test
  void logs_properties_and_redacts_sensitive_values(CapturedOutput output) throws Exception {
    StandardEnvironment environment = new StandardEnvironment();
    Map<String, Object> props = new HashMap<>();
    props.put("property-logger.enabled", "true");
    props.put("sample.secret", "should-not-log");
    props.put("sample.value", "visible");
    environment
        .getPropertySources()
        .addFirst(
            new MapPropertySource(
                "Config resource 'class path resource [application.yml]'", props));

    ApplicationPropertiesLogger logger = new ApplicationPropertiesLogger(environment);
    ReflectionTestUtils.setField(logger, "enabled", true);

    logger.run(new DefaultApplicationArguments(new String[0]));

    String logs = output.getOut();
    assertThat(logs).contains("===== Application YAML Properties at Startup =====");
    assertThat(logs).contains("property-logger.enabled = true");
    assertThat(logs).contains("sample.value = visible");
    assertThat(logs).contains("sample.secret = [!!!REDACTED!!!]");
  }
}
