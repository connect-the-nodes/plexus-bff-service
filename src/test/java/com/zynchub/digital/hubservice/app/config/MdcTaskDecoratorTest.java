package com.zynchub.digital.hubservice.app.config;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.concurrent.atomic.AtomicBoolean;
import org.junit.jupiter.api.Test;

class MdcTaskDecoratorTest {

  @Test
  void decorate_wraps_and_runs_runnable() {
    MdcTaskDecorator decorator = new MdcTaskDecorator();
    AtomicBoolean ran = new AtomicBoolean(false);

    Runnable decorated = decorator.decorate(() -> ran.set(true));
    decorated.run();

    assertThat(ran.get()).isTrue();
  }
}
