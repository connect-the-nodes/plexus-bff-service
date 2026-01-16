package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.concurrent.Executor;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.runner.ApplicationContextRunner;
import org.springframework.scheduling.concurrent.ThreadPoolTaskExecutor;
import org.springframework.test.util.ReflectionTestUtils;

class AsyncConfigurationTest {

  private final ApplicationContextRunner contextRunner =
      new ApplicationContextRunner().withUserConfiguration(AsyncConfiguration.class);

  @Test
  void task_executor_is_configured_with_mdc_decorator() {
    contextRunner.run(
        context -> {
          Executor executor = context.getBean("taskExecutor", Executor.class);
          assertThat(executor).isInstanceOf(ThreadPoolTaskExecutor.class);
          ThreadPoolTaskExecutor taskExecutor = (ThreadPoolTaskExecutor) executor;
          Object decorator = ReflectionTestUtils.getField(taskExecutor, "taskDecorator");
          assertThat(decorator).isInstanceOf(MdcTaskDecorator.class);
        });
  }
}
